package tlsserver

import "io/ioutil"

var S_CERT = `-----BEGIN CERTIFICATE-----
MIIDHjCCAoegAwIBAgIJAJH05SI8aNd0MA0GCSqGSIb3DQEBCwUAMIGnMQswCQYD
VQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTESMBAGA1UEBwwJU2FuIERpZWdv
MSAwHgYDVQQKDBdCbHVlIFNreSB3aGl0ZSBTYW5kIExMQzELMAkGA1UECwwCSVQx
HjAcBgNVBAMMFWJsdWUtc2t5LndoaXRlc2FuZC5pbzEgMB4GCSqGSIb3DQEJARYR
YWphcnZAZXhhbXBsZS5jb20wHhcNMTkwMjI4MDEzMTQxWhcNMjQwMjI3MDEzMTQx
WjCBpzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCkNhbGlmb3JuaWExEjAQBgNVBAcM
CVNhbiBEaWVnbzEgMB4GA1UECgwXQmx1ZSBTa3kgd2hpdGUgU2FuZCBMTEMxCzAJ
BgNVBAsMAklUMR4wHAYDVQQDDBVibHVlLXNreS53aGl0ZXNhbmQuaW8xIDAeBgkq
hkiG9w0BCQEWEWFqYXJ2QGV4YW1wbGUuY29tMIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQDwELpzGB7XElUH1cxJ8J7gOfRXoVDFu6hQ1f/YuHv1zAQk/zzCHoxy
VxU/ogtTEt7jqGYvHrIh4tjUwb2BNScAmWsPWkWlVBmG+pq1YenNlqOJ52tIZcZ5
YAu4V9N0Nz2nC8+ij7FrcuxbMEAFOVCf2JTy3iYatMm7b0ViPKmUvQIDAQABo1Aw
TjAdBgNVHQ4EFgQUg+vv2w007ilS0M3gDgyQ+0XvJCMwHwYDVR0jBBgwFoAUg+vv
2w007ilS0M3gDgyQ+0XvJCMwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOB
gQDgcGlcqIvRTQRyqgCgrurEG5KUy3tBlbqX6bnPsfMxqYcoWzIzLlrMWe8Hzlta
xJETxc+0v7wdj011z5anfxXonbTS5AF0NWFbtkC378tJ4Nm0N8FLO/+KV9iVtGKB
Ud84CYDn7LLKC5NnOm4MDmUCAdMKYqYNRhfmPSxxvLbStw==
-----END CERTIFICATE-----`

var S_KEY = `-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDwELpzGB7XElUH1cxJ8J7gOfRXoVDFu6hQ1f/YuHv1zAQk/zzC
HoxyVxU/ogtTEt7jqGYvHrIh4tjUwb2BNScAmWsPWkWlVBmG+pq1YenNlqOJ52tI
ZcZ5YAu4V9N0Nz2nC8+ij7FrcuxbMEAFOVCf2JTy3iYatMm7b0ViPKmUvQIDAQAB
AoGBANn4RmJUR0Q+R+haTifgi1DKLjoWpVE0ByqGc8viDeNqf2TcPt1+gUUcHpXt
Wtzt6GTKtSUZeOHdp8TduGQFz8cwVk04No6ygYVbV30VA9OYuYba0XxCyR0cKYf8
i2mV3pAfpobndec70wXTw6tn3GEIO/RISFY1cFWcu+scZyh9AkEA+Lmx/ooWD020
gDPOcXuzlQCIWoP/U+d40Cgi795YdpdDk2FDarkM5WGZnww4sym80HzAUn56kaS1
p4F60nBxcwJBAPcWMSbD8wFjlfOeyiWfVB0Ak1YREVI+IcpvMCQ7LD3TL914BxIW
vnDLO+0DOPZEYOxOkxOMeg9RY/nPW2MJlQ8CQQC+5uAL6vZlhpGcuKaiCXzbR05g
kuFdB9N9iODP1It3ckAmlUeGWUPhptie72Vxdf560tVWO8ddk9rtFv8rF6yrAkEA
oDmx0dOLR0FOwdYce90f7FatNEiJFO3Zd642Z6g/fi/ugA0PeLlq8TW5PG60h227
9EDXuvuDQ1+iFyJRvp0+HQJAHeireIT7yU+KY3GCdn2Vq4EjrxTu1MKn6loInkwo
JpF8fprPVvHTLj3c0+srGwlKgDzsOrNIh86D+nmA5ES/qw==
-----END RSA PRIVATE KEY-----`

func GetDefaultKeyFilePath() (filename string, err error) {
	f, err := ioutil.TempFile("/tmp", "key")
	if err != nil {
		return
	}
	filename = f.Name()
	_, err = f.Write([]byte(S_KEY))
	if err != nil {
		return
	}
	return
}
func GetDefaultCertFilePath() (filename string, err error) {
	f, err := ioutil.TempFile("/tmp", "cert")
	if err != nil {
		return
	}
	filename = f.Name()
	_, err = f.Write([]byte(S_CERT))
	if err != nil {
		return
	}
	return
}
