package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
)

// --- KULLANICI AYARLARI (BURALARI DOLDUR) ---
const MY_PRIVATE_KEY = "" // Örn: ac0974...
const KONTRAT_ADRESI = ""        // Örn: 0xe289...
const ANVIL_URL      = "http://127.0.0.1:8545"

// --- BITCOIN FONKSİYONLARI ---
// --- BITCOIN FONKSİYONLARI (DÜZELTİLMİŞ VERSİYON) ---
type BlockBilgisi struct {
	ID        string `json:"id"`
	Yukseklik int    `json:"height"`
}

func BitcoinSonBloguGetir() (*BlockBilgisi, error) {
	// API adresini değiştirdik: Son blokların listesini alıyoruz
	resp, err := http.Get("https://mempool.space/api/v1/blocks")
	if err != nil { 
		// Hatanın sebebini görmek için ekrana yazdıralım
		fmt.Println("İnternet Hatası:", err) 
		return nil, err 
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	
	// Gelen veri bir LİSTE (Array) olduğu için köşeli parantezli yapıyoruz
	var bloklar []BlockBilgisi
	err = json.Unmarshal(body, &bloklar)
	
	if err != nil {
		fmt.Println("JSON Çözme Hatası:", err)
		return nil, err
	}

	// Listenin ilk elemanı (en son blok) var mı?
	if len(bloklar) > 0 {
		return &bloklar[0], nil
	}
	
	return nil, fmt.Errorf("Veri listesi boş geldi")
}

// --- YARDIMCI: Fonksiyon İmzasını Bul ---
func GetMethodID(methodName string) []byte {
	transferFnSignature := []byte(methodName)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	return methodID
}

func main() {
	fmt.Println("🚀 RELAYER V4: HAZIRLANIYOR...")

	// 1. CÜZDAN KURULUMU
	privateKey, err := crypto.HexToECDSA(MY_PRIVATE_KEY)
	if err != nil { log.Fatal("❌ Private Key Hatası (başında 0x olmasın):", err) }

	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("👤 Bot Cüzdanı: %s\n", fromAddress.Hex())

	// 2. AĞA BAĞLAN
	client, err := ethclient.Dial(ANVIL_URL)
	if err != nil { log.Fatal("❌ Anvil'e bağlanılamadı. Terminalde 'anvil' çalışıyor mu?", err) }
	
	chainID, err := client.NetworkID(context.Background())
	if err != nil { log.Fatal(err) }
	fmt.Printf("🔗 Ağ ID: %s | Hedef Kontrat: %s\n", chainID, KONTRAT_ADRESI)

	// 3. SONSUZ DÖNGÜ
	var sonIslenenYukseklik int = 0

	for {
		// A. Bitcoin'e Bak
		btcBlok, err := BitcoinSonBloguGetir()
		if err != nil {
			fmt.Println("Bitcoin bekleniyor...")
			time.Sleep(5 * time.Second)
			continue
		}

		// B. Eğer yeni blok geldiyse
		if btcBlok.Yukseklik > sonIslenenYukseklik {
			fmt.Printf("\n📦 YENİ BITCOIN BLOK: %d\n", btcBlok.Yukseklik)
			fmt.Printf("   Hash: %s\n", btcBlok.ID)

			// --- İŞLEM GÖNDERME ---
			nonce, _ := client.PendingNonceAt(context.Background(), fromAddress)

			// Veriyi hazırla
			// YENİSİ:
            methodID := GetMethodID("blokGeldiParaBas(bytes32)")
			hashBytes, _ := hex.DecodeString(btcBlok.ID)
			
			var data []byte
			data = append(data, methodID...)
			data = append(data, hashBytes...)

			// Transaction Oluştur
			tx := types.NewTransaction(
				nonce,
				common.HexToAddress(KONTRAT_ADRESI),
				big.NewInt(0), 
				100000,        
				big.NewInt(1000000000), 
				data,          
			)

			// İmzala
			signedTx, _ := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)

			// Gönder
			fmt.Println("📤 Ethereum'a gönderiliyor...")
			err = client.SendTransaction(context.Background(), signedTx)
			
			if err != nil {
				fmt.Println("❌ GÖNDERME HATASI:", err)
			} else {
				fmt.Printf("✅ BAŞARILI! Tx Hash: %s\n", signedTx.Hash().Hex())
				sonIslenenYukseklik = btcBlok.Yukseklik
			}

		} else {
			fmt.Print(".") 
		}

		time.Sleep(3 * time.Second)
	}
}
