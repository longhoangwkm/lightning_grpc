package main

import (
    "io/ioutil"
    "fmt"
    "context"

    "lightning_grpc/lnrpc/github.com/lightningnetwork/lnd/lnrpc"
    "github.com/lightningnetwork/lnd/macaroons"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "gopkg.in/macaroon.v2"
)

func main() {
    tlsCertPath  := "/Users/tobi/Library/Application Support/Lnd_A/tls.cert"
    macaroonPath := "/Users/tobi/Library/Application Support/Lnd_A/data/chain/bitcoin/regtest/admin.macaroon"

    tlsCreds, err := credentials.NewClientTLSFromFile(tlsCertPath, "")
    if err != nil {
        fmt.Println("Cannot get node tls credentials", err)
        return
    }

    macaroonBytes, err := ioutil.ReadFile(macaroonPath)
    if err != nil {
        fmt.Println("Cannot read macaroon file", err)
        return
    }

    mac := &macaroon.Macaroon{}
    if err = mac.UnmarshalBinary(macaroonBytes); err != nil {
        fmt.Println("Cannot unmarshal macaroon", err)
        return
    }

    opts := []grpc.DialOption{
        grpc.WithTransportCredentials(tlsCreds),
        grpc.WithBlock(),
        grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(mac)),
    }

    conn, err := grpc.Dial("localhost:10009", opts...)
    if err != nil {
        fmt.Println("cannot dial to lnd", err)
        return
    }

    // client := lnrpc.NewLightningClient(conn)
    // ctx := context.Background()

    // getInfoResp, err := client.GetInfo(ctx, &lnrpc.GetInfoRequest{})
    // if err != nil {
    //     fmt.Println("Cannot get info from node:", err)
    //     return
    // }

    // fmt.Println(getInfoResp.BlockHeight)

    walletunlocker := lnrpc.NewWalletUnlockerClient(conn)
    ctx := context.Background()
    _, err = walletunlocker.UnlockWallet(ctx, &lnrpc.UnlockWalletRequest{WalletPassword: []byte("12345678")})
    if err != nil {
        fmt.Println(err.errorType())
        fmt.Println("Cannot unlock lightning node:", err)
        return
    } else {
        fmt.Println("Wallet is unlocked!")
    }
}
