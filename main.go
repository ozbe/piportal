package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			var request = struct {
				Active bool `json:"active"`
			}{}
			if err := json.Unmarshal(body, &request); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if request.Active {
				err = startVPN()
			} else {
				err = killVPN()
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
		default:
			j, err := json.Marshal(struct {
				Active bool `json:"is_active"`
			}{
				Active: isRunning(),
			})
			if err != nil {
				log.Fatal(err)
			}

			w.Write(j)
		}
	})
	addr := ":8081"
	log.Println("Listening at", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

var vpn *os.Process = nil

func isRunning() bool {
	return vpn != nil && vpn.Signal(syscall.Signal(0)) == nil
}

func killVPN() error {
	if !isRunning() {
		return nil
	}

	err := vpn.Signal(syscall.SIGTERM)
	_, _ = vpn.Wait()
	vpn = nil
	return err
}

func startVPN() error {
	if isRunning() {
		return nil
	}

	cmd := exec.Command("openvpn", "--config", "/etc/openvpn/us8605.nordvpn.com.udp.ovpn")
	cmd.Dir = "/etc/openvpn"
	if err := cmd.Start(); err != nil {
		return err
	}

	vpn = cmd.Process
	return nil
}
