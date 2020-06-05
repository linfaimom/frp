package plugin

import (
	"github.com/fatedier/frp/test/e2e/framework"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("[Feature: Client-Plugins]", func() {
	f := framework.NewDefaultFramework()

	Describe("UnixDomainSocket", func() {
		It("Expose a unix domain socket echo server", func() {
			localPortName := framework.TCPEchoServerPort
			serverConf := consts.DefaultServerConfig
			clientConf := consts.DefaultClientConfig

			getProxyConf := func(proxyName string, portName string, extra string) string {
				return fmt.Sprintf(`
				[%s]
				type = tcp
				local_port = {{ .%s }}
				remote_port = {{ .%s }}
				`+extra, proxyName, localPortName, portName)
			}

			tests := []struct {
				proxyName   string
				portName    string
				extraConfig string
			}{
				{
					proxyName: "tcp",
					portName:  framework.GenPortName("TCP"),
				},
				{
					proxyName:   "tcp-with-encryption",
					portName:    framework.GenPortName("TCPWithEncryption"),
					extraConfig: "use_encryption = true",
				},
				{
					proxyName:   "tcp-with-compression",
					portName:    framework.GenPortName("TCPWithCompression"),
					extraConfig: "use_compression = true",
				},
				{
					proxyName: "tcp-with-encryption-and-compression",
					portName:  framework.GenPortName("TCPWithEncryptionAndCompression"),
					extraConfig: `
					use_encryption = true
					use_compression = true
					`,
				},
			}

			// build all client config
			for _, test := range tests {
				clientConf += getProxyConf(test.proxyName, test.portName, test.extraConfig) + "\n"
			}
			// run frps and frpc
			f.RunProcesses([]string{serverConf}, []string{clientConf})

			for _, test := range tests {
				framework.ExpectTCPRequest(f.UsedPorts[test.portName],
					[]byte(consts.TestString), []byte(consts.TestString),
					connTimeout, test.proxyName)
			}

		})
	})
}
