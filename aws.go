package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type affiliate struct {
	public_key    string
	private_key   string
	associate_tag string
}

const (
	method  = "GET"
	uri     = "/onca/xml"
	host    = "webservices.amazon.com"
	version = "2011-08-01"
)

var (
	tags = []affiliate{
		{"AKIAJETPFP6RVOIBYNLQ", "JHi4SBbcoWkmCeJtJv234srD+OgoBFEb8qBcpIzt", "nosbl07-20"},
	}
)

func computeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func aws_signed_request(params map[string]string, public_key, private_key, associate_tag string) string {
	params["Service"] = "AWSECommerceService"
	params["AWSAccessKeyId"] = public_key

	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	//params["Timestamp"] = time.Unix(1442443227, 0).Format("2006-01-02T15:04:05Z")

	params["Version"] = version

	params["AssociateTag"] = associate_tag

	sorted_keys := make([]string, 0)
	for k, _ := range params {
		sorted_keys = append(sorted_keys, k)
	}

	// sort 'string' key in increasing order
	sort.Strings(sorted_keys)

	canonicalized_querys := make([]string, 0, 10)
	for _, key := range sorted_keys {
		key = strings.Replace(url.QueryEscape(key), "%7E", "~", -1)
		value := strings.Replace(url.QueryEscape(params[key]), "%7E", "~", -1)
		canonicalized_querys = append(canonicalized_querys, fmt.Sprintf("%s=%s", key, value))
	}

	canonicalized_query := strings.Join(canonicalized_querys, "&")

	// create the string to sign
	string_to_sign := method + "\n" + host + "\n" + uri + "\n" + canonicalized_query

	// calculate HMAC with SHA256 and base64-encoding
	signature := computeHmac256(string_to_sign, private_key)

	// encode the signature for the request
	signature = strings.Replace(url.QueryEscape(signature), "%7E", "~", -1)
	//signature = strings.Replace(signature, "%2F", "/", -1)

	return "http://" + host + uri + "?" + canonicalized_query + "&Signature=" + signature
}

func fetchAsinsInfo(index int) {
	//tag := tags[index]

}

type Result struct {
	Items ItemsS
}

type ItemsS struct {
	Request struct {
		ItemLookupRequest struct {
			ItemId []string
		}
	}
	Item []ItemS
}

type ItemS struct {
	ASIN       string
	ParentASIN string
	LargeImage struct {
		URL string
	}
	ImageSets struct {
		ImageSet []struct {
			LargeImage struct {
				URL string
			}
		}
	}
	ItemAttributes struct {
		Title             string
		Brand             string
		ProductGroup      string
		PackageDimensions struct {
			Weight int32
		}
		ListPrice struct {
			Amount int32
		}
	}
	Offers struct {
		Offer []OfferS
	}
	Variations VariationsS `xml:",omitempty"`
}

type VariationsS struct {
	Item []VariationItemS
}

type VariationItemS struct {
	ASIN           string
	ItemAttributes struct {
		Title     string
		Size      string
		Color     string
		ListPrice struct {
			Amount int32
		}
	}
	VariationAttributes struct {
		VariationAttribute []VariationAttributeS
	}
	Offers struct {
		Offer []OfferS
	}
}

type VariationAttributeS struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type OfferS struct {
	Merchant struct {
		Name string
	}
	OfferListing struct {
		Price struct {
			Amount int
		}
		IsEligibleForSuperSaverShipping bool
	}
}

type AwsInfo struct {
	ASIN         string  `json:"a"`
	ParentASIN   string  `json:"h,omitempty"`
	LargeImage   string  `json:"l,omitempty"`
	Title        string  `json:"e,omitempty"`
	Brand        string  `json:"b,omitempty"`
	ProductGroup string  `json:"g,omitempty"`
	Weight       float32 `json:"w,omitempty"`
	ListPrice    float32 `json:"o,omitempty"`
	Price        float32 `json:"n,omitempty"`

	Variations map[string]SubAwsInfo `json:"j,omitempty"`
}

type SubAwsInfo struct {
	Title      string  `json:"e,omitempty"`
	LargeImage string  `json:"l,omitempty"`
	ListPrice  float32 `json:"o,omitempty"`
	Price      float32 `json:"n,omitempty"`
	Merchant   int32   `json:"k,omitempty"`

	VariationAttributes []VariationAttributeS `json:"attrs,omitempty"`
}

func parseAwsInfo(content []byte) (interface{}, error) {
	var result Result
	err := xml.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(content))
	fmt.Println(result)
	detail := make(map[string]*AwsInfo)
	for _, item := range result.Items.Item {
		asin := item.ASIN
		detail[asin] = &AwsInfo{}
		detail[asin].ASIN = asin
		if item.ParentASIN != "" {
			detail[asin].ParentASIN = item.ParentASIN
		}
		if item.LargeImage.URL != "" {
			detail[asin].LargeImage = item.LargeImage.URL
		}
		if item.ItemAttributes.Title != "" {
			detail[asin].Title = item.ItemAttributes.Title
		}
		if item.ItemAttributes.Brand != "" {
			detail[asin].Brand = item.ItemAttributes.Brand
		}
		if item.ItemAttributes.ProductGroup != "" {
			detail[asin].LargeImage = item.ItemAttributes.ProductGroup
		}
		if item.ItemAttributes.PackageDimensions.Weight != 0 {
			detail[asin].Weight = float32(item.ItemAttributes.PackageDimensions.Weight)
		}
		if item.ItemAttributes.ListPrice.Amount != 0 {
			detail[asin].ListPrice = float32(item.ItemAttributes.ListPrice.Amount)
		}

	}
	return detail, nil
}

func ItemLookup(asins string, responseGroup string, aff *affiliate) (interface{}, error) {
	params := make(map[string]string)
	params["Operation"] = "ItemLookup"
	params["ItemId"] = asins
	params["ResponseGroup"] = responseGroup
	request := aws_signed_request(params, aff.public_key, aff.private_key, aff.associate_tag)
	resp, err := http.Get(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseAwsInfo(body)
}
