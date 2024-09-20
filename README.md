# Hyperledger Fabric MeetCoin v2

### Cara pakai:

1. Install prereq HLF sesuai OS masing-masing, bisa dibaca [di sini](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)

2. Buka terminal, install HLF dengan command

```
./install-fabric.sh docker binary
```

3. Setelah selesai, pindah ke direktori `rest-api-go` dan install dependency-nya

4. Balik ke root folder, pindah ke direktori ke `test-network` dan ketik command

```
./startNetwork.sh
```

5. Pindah ke direktori `test-network/fabric-gateway-api-go` dan ketik command untuk start Fabric Gateway API

```
go run main.go
```

6. Pindah lagi ke direktori `rest-api-go` dan ketik command untuk start REST server

```
go run main.go
```

7. REST server akan up di address

```
http://localhost:8080
```

8. Untuk terminate REST server & Fabric Gateway API, bisa `Ctrl+C` di terminal di direktori `rest-api-go` & `test-network/fabric-gateway-api-go`

9. Untuk stop jaringan HLF, pindah ke direktori `test-network` dan ketik command `./network.sh down`

#### Note:

Kalau misal ada error permission denied, coba ubah permission file di direktori root, `test-network`, `test-network/scripts`, dan `test-network/organizations` dengan command

```
chmod 755 *.sh
```

Pastikan dockernya berjalan sebelum start jaringan HLF

Pastikan juga start jaringan & running REST server pakai terminal WSL2 (khusus Windows)
