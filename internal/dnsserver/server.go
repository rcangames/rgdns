package dnsserver
import("encoding/json";"fmt";"net";"os";"strings";"time";"github.com/miekg/dns")
type configFile struct{Zones []zone `json:"zones"`}
type zone struct{Name string `json:"name"`;TTL uint32 `json:"ttl"`;Records []record `json:"records"`}
type record struct{Name string `json:"name"`;Type string `json:"type"`;Value string `json:"value"`}
type Listener struct{server *dns.Server}
type Server struct{zones map[string]zone}
func New(path string)(*Server,error){b,e:=os.ReadFile(path);if e!=nil{return nil,e};var c configFile;if e=json.Unmarshal(b,&c);e!=nil{return nil,e};s:=&Server{zones:map[string]zone{}};for _,z:=range c.Zones{z.Name=dns.Fqdn(strings.TrimSpace(z.Name));if z.TTL==0{z.TTL=300};for i:=range z.Records{r:=&z.Records[i];r.Name=dns.Fqdn(strings.TrimSpace(r.Name));r.Type=strings.ToUpper(strings.TrimSpace(r.Type));if !dns.IsSubDomain(z.Name,r.Name){return nil,fmt.Errorf("record %s outside zone %s",r.Name,z.Name)}};s.zones[strings.ToLower(z.Name)]=z};if len(s.zones)==0{return nil,fmt.Errorf("no zones configured")};return s,nil}
func(s *Server)Listen(addr,network string)*Listener{return &Listener{server:&dns.Server{Addr:addr,Net:network,Handler:s}}}
func(l *Listener)Activate()error{return l.server.ListenAndServe()}
func(l *Listener)Shutdown()error{return l.server.Shutdown()}
func(s *Server)ServeDNS(w dns.ResponseWriter,req *dns.Msg){m:=new(dns.Msg);m.SetReply(req);m.Authoritative=true;if len(req.Question)==0{m.Rcode=dns.RcodeFormatError;_=w.WriteMsg(m);return};q:=req.Question[0];z,ok:=s.findZone(q.Name);if !ok{m.Rcode=dns.RcodeRefused;_=w.WriteMsg(m);return};found:=false;for _,r:=range z.Records{if !strings.EqualFold(r.Name,q.Name){continue};found=true;rr,e:=toRR(r,z.TTL);if e!=nil{m.Rcode=dns.RcodeServerFailure;_=w.WriteMsg(m);return};if q.Qtype==dns.TypeANY||rr.Header().Rrtype==q.Qtype{m.Answer=append(m.Answer,rr)}};if !found||len(m.Answer)==0{m.Ns=append(m.Ns,soa(z))};if !found{m.Rcode=dns.RcodeNameError};_=w.WriteMsg(m)}
func(s *Server)findZone(name string)(zone,bool){name=strings.ToLower(dns.Fqdn(name));var best zone;n:=-1;for k,z:=range s.zones{if dns.IsSubDomain(k,name)&&len(k)>n{best,n=z,len(k)}};return best,n>=0}
func toRR(r record,ttl uint32)(dns.RR,error){h:=fmt.Sprintf("%s %d IN %s",dns.Fqdn(r.Name),ttl,r.Type);switch r.Type{case "A":ip:=net.ParseIP(r.Value).To4();if ip==nil{return nil,fmt.Errorf("invalid IPv4")};return dns.NewRR(h+" "+ip.String());case "AAAA":ip:=net.ParseIP(r.Value).To16();if ip==nil||ip.To4()!=nil{return nil,fmt.Errorf("invalid IPv6")};return dns.NewRR(h+" "+ip.String());case "CNAME","NS":return dns.NewRR(h+" "+dns.Fqdn(r.Value));case "TXT":return dns.NewRR(h+fmt.Sprintf(" %q",r.Value));default:return nil,fmt.Errorf("unsupported type %s",r.Type)}}
func soa(z zone)dns.RR{return &dns.SOA{Hdr:dns.RR_Header{Name:z.Name,Rrtype:dns.TypeSOA,Class:dns.ClassINET,Ttl:z.TTL},Ns:"ns1."+z.Name,Mbox:"hostmaster."+z.Name,Serial:uint32(time.Now().Unix()),Refresh:3600,Retry:600,Expire:86400,Minttl:z.TTL}}
