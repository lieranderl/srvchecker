package output


func MakeTcpConnectivity(discoveredsrv []Srv, channel chan []Fqdn){
	fqdntcpconnectivity_list := make([]Fqdn,0)
	fqdn_list := make([]string, 0)
	for _, srv := range discoveredsrv {
		for _, fqdn := range srv.Fqdns {
			if !stringInSlice(fqdn.Name, fqdn_list) {
				fqdntcpconnectivity_list = append(fqdntcpconnectivity_list, fqdn)
			}
			fqdn_list = append(fqdn_list, fqdn.Name)
		}
	}

	channel <- fqdntcpconnectivity_list
}