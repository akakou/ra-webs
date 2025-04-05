package monitor

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"fmt"

	"github.com/akakou/ra-webs/monitor/ent"
	"github.com/akakou/ra-webs/monitor/ent/taserver"
)

type publicKey interface{}

func (monitor *Monitor) Check(pk publicKey, id int) {
	serv, err := monitor.DB.Client.TAServer.
		Query().
		Order(ent.Desc(taserver.FieldID)).
		First(*monitor.DB.Ctx)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	serv.Update().SetMonitorLogID(id).Save(*monitor.DB.Ctx)

	unmarshaledPublicKey, isRSA := pk.(*rsa.PublicKey)
	if !isRSA {
		fmt.Printf("Violation: %v\n", errPublicKeyNotRSA)
		monitor.Revoke(serv)
		return
	}

	publicKey := x509.MarshalPKCS1PublicKey(unmarshaledPublicKey)
	fmt.Printf("compireing public key:\n%v\n!=%v\n\n", serv.PublicKey, publicKey)

	if !bytes.Equal(serv.PublicKey, publicKey) {
		fmt.Printf("Violation: %v: %v != %v", errPublicKeyNotMatch, serv.PublicKey, publicKey)
		monitor.Revoke(serv)
		return
	}

	if !serv.IsActive {
		serv.Update().SetIsActive(true).Save(*monitor.DB.Ctx)
	}
}
