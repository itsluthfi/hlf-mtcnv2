# MeetCoin v2

### Cara pakai:

1. Install prereq HLF sesuai OS masing-masing, bisa dibaca [di sini](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html)

2. Buka terminal, install HLF dengan command

```
./install-fabric.sh docker binary
```

3. Setelah selesai, pindah ke direktori `mtcn-chaincode-go` & `mtcn-rest-api-go` dan install dependency-nya

4. Balik ke root folder, pindah ke direktori ke `test-network` dan ketik command

```
./startNetwork.sh
```

5. Setelah selesai, pindah lagi ke direktori `mtcn-rest-api-go` dan ketik command untuk start REST server

```
go run main.go
```

6. Setelah server berhasil up, ada dua endpoint yang bisa diakses

```
http://localhost:3000/invoke
atau
http://localhost:3000/query
```

7. Untuk terminate server, bisa `Ctrl+C` di terminal di direktori `mtcn-rest-api-go`

8. Untuk stop jaringan HLF, pindah ke direktori `test-network` dan ketik command `./network.sh down`

#### Note:

Kalau misal ada error permission denied, coba ubah permission file di direktori root, `test-network`, `test-network/scripts`, dan `test-network/organizations` dengan command

```
chmod 755 *.sh
```
