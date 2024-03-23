package zota

import "testing"

func setupZotaDepositRequest(merchantOrderId, orderAmount, customerEmail string) ZotaDepositRequest {
    return ZotaDepositRequest {
        MerchantOrderID: merchantOrderId,
        MerchantOrderDesc: "Test deposit",
        OrderAmount: orderAmount,

        CustomerEmail: customerEmail,
		CustomerFirstName: "Nikola",
		CustomerLastName:  "Velichkov",
		CustomerIP:        "127.0.0.1",
		CustomerPhone:     "+4511111111",

		CustomerAddress:     "Line",
		CustomerCountryCode: "DK",
		CustomerCity:        "City",
		CustomerZipCode:     "zip",

        RedirectUrl: "test",
        CheckoutUrl: "test",
        Signature: "",
    }
}


type depositGenSignatureTest struct {
    endpointId string
    merchantOrderId string
    orderAmount string
    customerEmail string
    merchantSecretKey string
    expected string
}

// Expected values have been generated from https://mg-tools.zotapay.com
var depositGenSignatureTests = []depositGenSignatureTest {
    {
        endpointId: "111111",
        merchantOrderId: "e31edd0d-76a6-4f1c-be19-4504ff5b89d7",
        orderAmount: "13.37",
        customerEmail: "federlizer@protonmail.com",
        merchantSecretKey: "00000000-1111-2222-3333-444444444444",
        expected: "8853fd7c21545ede527ee3e37f2c384b7d2cb4bc2b027c25c472b1fc38e4608e",
    },

    {
        endpointId: "234567",
        merchantOrderId: "e31edd0d-76a6-4f1c-be19-4504ff5b89d7",
        orderAmount: "1300.37",
        customerEmail: "test@gmail.com",
        merchantSecretKey: "00000000-1111-2222-3333-444444444444",
        expected: "e869a8ccb65bd247bbd1cd1935ac960a86ff2fa0a1450f1abb17180b38b7bde9",
    },
}

func TestDepositGenSignature(t *testing.T) {
    for _, test := range depositGenSignatureTests {
        request := setupZotaDepositRequest(test.merchantOrderId, test.orderAmount, test.customerEmail)
        signature := request.GenSignature(test.endpointId, test.merchantSecretKey)

        if signature != test.expected {
            t.Errorf("Output %q does not equal expected %q\n", signature, test.expected)
        }
    }
}
