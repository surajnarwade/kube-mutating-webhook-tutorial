package main

import (
	"crypto/tls"
	"log"
	"sync"

	"gopkg.in/fsnotify.v1"
)

type keypairReloader struct {
	certMu   sync.RWMutex
	cert     *tls.Certificate
	certPath string
	keyPath  string
}

func NewKeypairReloader(certPath, keyPath string) (*keypairReloader, error) {
	result := &keypairReloader{
		certPath: certPath,
		keyPath:  keyPath,
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}
	result.cert = &cert

	// pathChange := make(chan bool)

	// err = watchFiles(pathChange)
	// if err != nil {
	// 	return nil, err
	// }
	// creates a new file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	// defer watcher.Close()

	if err := watcher.Add("/etc/webhook/certs"); err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				log.Print(event)
				if event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("File Changed, reloading TLS certificate and key from and ")
					if err := result.maybeReload(); err != nil {
						log.Printf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
					}
					// pathChange <- true
				}
				// fmt.Printf("EVENT! %#v\n", event)

				// watch for errors
			case err := <-watcher.Errors:
				log.Print("error", err)
				// fmt.Println("ERROR", err)
			}
		}
	}()

	// go func() {
	// c := make(chan os.Signal, 1)
	// signal.Notify(c, syscall.SIGHUP)
	// for range c {
	// 	log.Printf("Received SIGHUP, reloading TLS certificate and key from and ")
	// 	if err := result.maybeReload(); err != nil {
	// 		log.Printf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
	// 	}
	// }
	// for {

	// 	// c := <-pathChange
	// 	if c {
	// 		log.Printf("File Changed, reloading TLS certificate and key from and ")
	// 		if err := result.maybeReload(); err != nil {
	// 			log.Printf("Keeping old TLS certificate because the new one could not be loaded: %v", err)
	// 		}
	// 		// pathChange <- false

	// 	}
	// }
	// }()
	return result, nil
}

func (kpr *keypairReloader) maybeReload() error {
	newCert, err := tls.LoadX509KeyPair(kpr.certPath, kpr.keyPath)
	if err != nil {
		return err
	}
	kpr.certMu.Lock()
	defer kpr.certMu.Unlock()
	kpr.cert = &newCert
	return nil
}

func (kpr *keypairReloader) GetCertificateFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		kpr.certMu.RLock()
		defer kpr.certMu.RUnlock()
		return kpr.cert, nil
	}
}

// func watchFiles(pathChange chan<- bool) error {
// 	// creates a new file watcher
// 	watcher, err := fsnotify.NewWatcher()
// 	if err != nil {
// 		return err
// 	}
// 	// defer watcher.Close()

// 	if err := watcher.Add("/etc/webhook/certs"); err != nil {
// 		return err
// 	}

// 	go func() {
// 		for {
// 			select {
// 			// watch for events
// 			case event := <-watcher.Events:
// 				if event.Op&fsnotify.Write == fsnotify.Write {
// 					log.Print(event)
// 					pathChange <- true
// 				}
// 				// fmt.Printf("EVENT! %#v\n", event)

// 				// watch for errors
// 			case err := <-watcher.Errors:
// 				log.Print("error", err)
// 				// fmt.Println("ERROR", err)
// 			}
// 		}
// 	}()

// 	return nil
// }
