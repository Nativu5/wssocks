package client

import (
	"errors"
	"flag"
	"github.com/genshen/cmds"
	"github.com/genshen/wssocks/wss"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

const CommandNameClient = "client"

var clientCommand = &cmds.Command{
	Name:        CommandNameClient,
	Summary:     "run as client mode",
	Description: "run as client program.",
	CustomFlags: false,
	HasOptions:  true,
}

func init() {
	var client client
	fs := flag.NewFlagSet(CommandNameClient, flag.ExitOnError)
	clientCommand.FlagSet = fs
	clientCommand.FlagSet.StringVar(&client.address, "addr", ":1080", `listen address of socks5.`)
	clientCommand.FlagSet.BoolVar(&client.http, "http", false, `enable http and https proxy.`)
	clientCommand.FlagSet.StringVar(&client.httpAddr, "http-addr", ":1086", `listen address of http and https(if enabled).`)
	clientCommand.FlagSet.StringVar(&client.remote, "remote", "", `server address and port(e.g: ws://example.com:1088).`)

	clientCommand.FlagSet.Usage = clientCommand.Usage // use default usage provided by cmds.Command.
	clientCommand.Runner = &client

	cmds.AllCommands = append(cmds.AllCommands, clientCommand)
}

type client struct {
	address   string   // local listening address
	http      bool     // enable http and https proxy
	httpAddr  string   // listen address of http and https(if it is enabled)
	remote    string   // string usr of server
	remoteUrl *url.URL // url of server
	//	remoteHeader http.Header
}

func (c *client) PreRun() error {
	// check remote address
	if c.remote == "" {
		return errors.New("empty remote address")
	}
	if u, err := url.Parse(c.remote); err != nil {
		return err
	} else {
		c.remoteUrl = u
	}

	if c.http {
		log.Info("http(s) proxy is enabled.")
	} else {
		log.Info("http(s) proxy is disabled.")
	}
	return nil
}

func (c *client) Run() error {
	// start websocket connection (to remote server).
	log.WithFields(log.Fields{
		"remote": c.remoteUrl.String(),
	}).Info("connecting to wssocks server.")

	dialer := websocket.DefaultDialer
	wsHeader := make(http.Header) // header in websocket request(default is nil)

	// loading and execute plugin
	if clientPlugin.HasPlugin() {
		// in the plugin, we may add http header/dialer and modify remote address.
		if err := clientPlugin.RedirectPlugin.BeforeRequest(dialer, c.remoteUrl, wsHeader); err != nil {
			return err
		}
	}

	wsc, err := wss.NewWebSocketClient(websocket.DefaultDialer, c.remoteUrl.String(), wsHeader)
	if err != nil {
		log.Fatal("establishing connection error:", err)
	}
	log.WithFields(log.Fields{
		"remote": c.remoteUrl.String(),
	}).Info("connected to wssocks server.")
	// todo chan for wsc and tcp accept
	defer wsc.WSClose()

	// negotiate version
	if version, err := wss.ExchangeVersion(wsc.WsConn); err != nil {
		return err
	} else {
		log.WithFields(log.Fields{
			"version code":   version.VersionCode,
			"version number": version.Version,
		}).Info("server version")

		if version.VersionCode != wss.VersionCode {
			return errors.New("incompatible protocol version of client and server")
		}
		if version.Version != wss.CoreVersion {
			log.WithFields(log.Fields{
				"client wssocks version": wss.CoreVersion,
				"server wssocks version": version.Version,
			}).Warning("different version of client and server wssocks")
		}
	}

	// start websocket message listen.
	go func() {
		if err := wsc.ListenIncomeMsg(); err != nil {
			log.Error("error websocket read:", err)
		}
	}()
	// send heart beats.
	go func() {
		if err := wsc.HeartBeat(); err != nil {
			log.Info("heartbeat ending", err)
		}
	}()

	// http listening
	if c.http {
		go func() {
			handle := wss.NewHttpProxy(wsc);
			if err := http.ListenAndServe(c.httpAddr, &handle); err != nil {
				log.Fatalln(err)
			}
		}()
	}

	// start listen for socks5 and https connection.
	if err := wss.ListenAndServe(wsc, c.address, c.http); err != nil {
		return err
	}
	return nil
}
