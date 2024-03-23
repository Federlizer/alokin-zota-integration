package zota

import "testing"

func setupZotaOrderStatusRequest(orderId, merchantOrderId string, timestamp int64) ZotaOrderStatusRequest {
    return ZotaOrderStatusRequest {
        OrderId: orderId,
        MerchantOrderId: merchantOrderId,
        Timestamp: timestamp,
    }
}

type orderStatusGenSignatureTest struct {
    merchantId string
    secretKey string
    orderId string
    merchantOrderId string
    timestamp int64
    expected string
}

// Expected values have been generated from https://mg-tools.zotapay.com
var orderStatusGenSignatureTests = []orderStatusGenSignatureTest {
     {
        merchantId: "BOOOGYBOO123",
        secretKey: "00000000-1111-2222-3333-444444444444",
        orderId: "e7d7f9d6c35005858b58d5455355411171c0c255",
        merchantOrderId: "59fa8d26-2a16-4665-963a-65fd5c0d9da2",
        timestamp: 1711199920,
        expected: "baa4fd0980da606168d201c34795ed8837f04404891d506848328e6426e77a04",
    },

    {
        merchantId: "MYMERCHANTID",
        secretKey: "00000000-1111-2222-3333-444444444444",
        orderId: "9697960f4561f634dd5363590e55c93586a3721e",
        merchantOrderId: "43590438-61d1-4e7a-a31b-6df4772d9b9a",
        timestamp: 1711199946,
        expected: "e21cf2ee283f7f6c01d250974061c2e40c454e6dbaf9c99ecf3d3c7ae9e2bd7e",
    },
}

func TestGenSignature(t *testing.T) {
    for _, test := range orderStatusGenSignatureTests {
        request := setupZotaOrderStatusRequest(test.orderId, test.merchantOrderId, test.timestamp)
        signature := request.GenSignature(test.merchantId, test.secretKey)

        if signature != test.expected {
            t.Errorf("Output %q does not equal expected %q\n", signature, test.expected)
        }
    }
}
