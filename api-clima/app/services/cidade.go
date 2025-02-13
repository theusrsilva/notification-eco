package services

import (
	"api-clima/framework/database"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type CidadesResponse struct {
	XMLName xml.Name `xml:"cidades"`
	Cidades []Cidade `xml:"cidade"`
}
type Cidade struct {
	Nome string `xml:"nome"`
	UF   string `xml:"uf"`
	ID   string `xml:"id"`
}
type CidadeResponse struct {
	ID   string `json:"id"`
	Nome string `json:"nome"`
	UF   string `json:"uf"`
}

func GetCidades(cidade string) ([]Cidade, error) {
	buscaCidade := strings.ReplaceAll(cidade, " ", "%20")
	url := fmt.Sprintf(os.Getenv("CPTEC_URL")+"listaCidades?city=%s", buscaCidade)

	redis := database.NewRedisClient(0)

	valoresCache, err := redis.Find(buscaCidade + time.Now().Format("2006-01-02"))

	if valoresCache != "" {
		var cidadesCache []Cidade
		err := json.Unmarshal([]byte(valoresCache), &cidadesCache)
		if err != nil {
			return nil, fmt.Errorf("erro ao desserializar dados do cache: %v", err)
		}
		return cidadesCache, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição: %v", err)
	}
	defer resp.Body.Close()

	reader, err := charset.NewReader(resp.Body, "ISO-8859-1")
	if err != nil {
		return nil, fmt.Errorf("erro ao configurar leitor de charset: %v", err)
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	decoder := xml.NewDecoder(strings.NewReader(string(body)))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		if charset == "ISO-8859-1" {
			return input, nil
		}
		return nil, fmt.Errorf("não suportado charset: %s", charset)
	}
	var cidadesResponse CidadesResponse
	err = decoder.Decode(&cidadesResponse)

	var cidadesEncontradas []Cidade
	for _, c := range cidadesResponse.Cidades {
		cidadesEncontradas = append(cidadesEncontradas, c)
	}
	if len(cidadesEncontradas) == 0 {
		return nil, fmt.Errorf("nenhuma cidade encontrada")
	}

	cidadesCacheJSON, err := json.Marshal(cidadesEncontradas)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar dados para cache: %v", err)
	}
	err = redis.Insert(buscaCidade+time.Now().Format("2006-01-02"), string(cidadesCacheJSON), 24*time.Hour)
	return cidadesEncontradas, nil

}
