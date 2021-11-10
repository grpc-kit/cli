// Copyright © 2020 Li MingQing <mingqing@henji.org>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

func (t *templateService) fileDirectoryRoot() {
	t.files = append(t.files, &templateFile{
		name: ".gitignore",
		body: `
# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories
vendor/

# Config file and certificate
config/*.pem
config/*.crt
config/*.key
config/app-dev-*
config/app-test-*
config/app-prod-*

# Others
.swp
.bak
.idea/
build/
backup/
scripts/env-dev-*.sh
scripts/env-test-*.sh
scripts/env-prod-*.sh
api/assets_vfsdata.go
api/doc/openapi-spec/*
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "go.mod",
		parse: true,
		body: `
module {{ .Global.GitDomain }}/{{ .Global.ProductCode }}/{{ .Global.ShortName }}

go 1.13

require (
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/grpc-kit/api v0.0.0-20200430142535-fa4f5e68cf77
	github.com/grpc-kit/cfg v0.1.6
	github.com/grpc-kit/pkg v0.1.1
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.38.0
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.5
	google.golang.org/genproto => google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
`,
	})

	t.files = append(t.files, &templateFile{
		name:  "go.sum",
		parse: true,
		body: `
cloud.google.com/go v0.34.0/go.mod h1:aQUYkXzVsufM+DwF1aE+0xfcU+56JwCaLick0ClmMTw=
cloud.google.com/go v0.38.0/go.mod h1:990N+gfupTy94rShfmMCWGDn0LpTmnzTp2qbd1dvSRU=
cloud.google.com/go v0.44.1/go.mod h1:iSa0KzasP4Uvy3f1mN/7PiObzGgflwredwwASm/v6AU=
cloud.google.com/go v0.44.2/go.mod h1:60680Gw3Yr4ikxnPRS/oxxkBccT6SA1yMk63TGekxKY=
cloud.google.com/go v0.45.1/go.mod h1:RpBamKRgapWJb87xiFSdk4g1CME7QZg3uwTez+TSTjc=
cloud.google.com/go v0.46.3/go.mod h1:a6bKKbmY7er1mI7TEI4lsAkts/mkhTSZK8w33B4RAg0=
cloud.google.com/go v0.50.0/go.mod h1:r9sluTvynVuxRIOHXQEHMFffphuXHOMZMycpNR5e6To=
cloud.google.com/go v0.52.0/go.mod h1:pXajvRH/6o3+F9jDHZWQ5PbGhn+o8w9qiu/CffaVdO4=
cloud.google.com/go v0.53.0/go.mod h1:fp/UouUEsRkN6ryDKNW/Upv/JBKnv6WDthjR6+vze6M=
cloud.google.com/go v0.54.0/go.mod h1:1rq2OEkV3YMf6n/9ZvGWI3GWw0VoqH/1x2nd8Is/bPc=
cloud.google.com/go v0.56.0/go.mod h1:jr7tqZxxKOVYizybht9+26Z/gUq7tiRzu+ACVAMbKVk=
cloud.google.com/go v0.57.0/go.mod h1:oXiQ6Rzq3RAkkY7N6t3TcE6jE+CIBBbA36lwQ1JyzZs=
cloud.google.com/go v0.62.0/go.mod h1:jmCYTdRCQuc1PHIIJ/maLInMho30T/Y0M4hTdTShOYc=
cloud.google.com/go v0.65.0/go.mod h1:O5N8zS7uWy9vkA9vayVHs65eM1ubvY4h553ofrNHObY=
cloud.google.com/go v0.72.0/go.mod h1:M+5Vjvlc2wnp6tjzE102Dw08nGShTscUx2nZMufOKPI=
cloud.google.com/go v0.74.0/go.mod h1:VV1xSbzvo+9QJOxLDaJfTjx5e+MePCpCWwvftOeQmWk=
cloud.google.com/go v0.78.0/go.mod h1:QjdrLG0uq+YwhjoVOLsS1t7TW8fs36kLs4XO5R5ECHg=
cloud.google.com/go v0.79.0/go.mod h1:3bzgcEeQlzbuEAYu4mrWhKqWjmpprinYgKJLgKHnbb8=
cloud.google.com/go v0.81.0 h1:at8Tk2zUz63cLPR0JPWm5vp77pEZmzxEQBEfRKn1VV8=
cloud.google.com/go v0.81.0/go.mod h1:mk/AM35KwGk/Nm2YSeZbxXdrNK3KZOYHmLkOqC2V6E0=
cloud.google.com/go/bigquery v1.0.1/go.mod h1:i/xbL2UlR5RvWAURpBYZTtm/cXjCha9lbfbpx4poX+o=
cloud.google.com/go/bigquery v1.3.0/go.mod h1:PjpwJnslEMmckchkHFfq+HTD2DmtT67aNFKH1/VBDHE=
cloud.google.com/go/bigquery v1.4.0/go.mod h1:S8dzgnTigyfTmLBfrtrhyYhwRxG72rYxvftPBK2Dvzc=
cloud.google.com/go/bigquery v1.5.0/go.mod h1:snEHRnqQbz117VIFhE8bmtwIDY80NLUZUMb4Nv6dBIg=
cloud.google.com/go/bigquery v1.7.0/go.mod h1://okPTzCYNXSlb24MZs83e2Do+h+VXtc4gLoIoXIAPc=
cloud.google.com/go/bigquery v1.8.0/go.mod h1:J5hqkt3O0uAFnINi6JXValWIb1v0goeZM77hZzJN/fQ=
cloud.google.com/go/datastore v1.0.0/go.mod h1:LXYbyblFSglQ5pkeyhO+Qmw7ukd3C+pD7TKLgZqpHYE=
cloud.google.com/go/datastore v1.1.0/go.mod h1:umbIZjpQpHh4hmRpGhH4tLFup+FVzqBi1b3c64qFpCk=
cloud.google.com/go/firestore v1.1.0/go.mod h1:ulACoGHTpvq5r8rxGJ4ddJZBZqakUQqClKRT5SZwBmk=
cloud.google.com/go/pubsub v1.0.1/go.mod h1:R0Gpsv3s54REJCy4fxDixWD93lHJMoZTyQ2kNxGRt3I=
cloud.google.com/go/pubsub v1.1.0/go.mod h1:EwwdRX2sKPjnvnqCa270oGRyludottCI76h+R3AArQw=
cloud.google.com/go/pubsub v1.2.0/go.mod h1:jhfEVHT8odbXTkndysNHCcx0awwzvfOlguIAii9o8iA=
cloud.google.com/go/pubsub v1.3.1/go.mod h1:i+ucay31+CNRpDW4Lu78I4xXG+O1r/MAHgjpRVR+TSU=
cloud.google.com/go/storage v1.0.0/go.mod h1:IhtSnM/ZTZV8YYJWCY8RULGVqBDmpoyjwiyrjsg+URw=
cloud.google.com/go/storage v1.5.0/go.mod h1:tpKbwo567HUNpVclU5sGELwQWBDZ8gh0ZeosJ0Rtdos=
cloud.google.com/go/storage v1.6.0/go.mod h1:N7U0C8pVQ/+NIKOBQyamJIeKQKkZ+mxpohlUTyfDhBk=
cloud.google.com/go/storage v1.8.0/go.mod h1:Wv1Oy7z6Yz3DshWRJFhqM/UCfaWIRTdp0RXyy7KQOVs=
cloud.google.com/go/storage v1.10.0/go.mod h1:FLPqc6j+Ki4BU591ie1oL6qBQGu2Bl/tZ9ullr3+Kg0=
dmitri.shuralyov.com/gpu/mtl v0.0.0-20190408044501-666a987793e9/go.mod h1:H6x//7gZCb22OMCxBHrMx7a5I7Hp++hsVxbQ4BYO7hU=
github.com/BurntSushi/toml v0.3.1/go.mod h1:xHWCNGjB5oqiDr8zfno3MHue2Ht5sIBksp03qcyfWMU=
github.com/BurntSushi/xgb v0.0.0-20160522181843-27f122750802/go.mod h1:IVnqGOEym/WlBOVXweHU+Q+/VP0lqqI8lqeDx9IjBqo=
github.com/NYTimes/gziphandler v0.0.0-20170623195520-56545f4a5d46/go.mod h1:3wb06e3pkSAbeQ52E9H9iFoQsEEwGN64994WTCIhntQ=
github.com/OneOfOne/xxhash v1.2.2/go.mod h1:HSdplMjZKSmBqAxg5vPj2TmRDmfkzw+cTzAElWljhcU=
github.com/PuerkitoBio/purell v1.1.1/go.mod h1:c11w/QuzBsJSee3cPx9rAFu61PvFxuPbtSwDGJws/X0=
github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578/go.mod h1:uGdkoq3SwY9Y+13GIhn11/XLaGBb4BfwItxLd5jeuXE=
github.com/alecthomas/template v0.0.0-20160405071501-a0175ee3bccc/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751/go.mod h1:LOuyumcjzFXgccqObfd/Ljyb9UuFJ6TxHnclSeseNhc=
github.com/alecthomas/units v0.0.0-20151022065526-2efee857e7cf/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
github.com/alecthomas/units v0.0.0-20190717042225-c3de453c63f4/go.mod h1:ybxpYRFXyAe+OPACYpWeL0wqObRcbAqCMya13uyzqw0=
github.com/antihax/optional v0.0.0-20180407024304-ca021399b1a6/go.mod h1:V8iCPQYkqmusNa815XgQio277wI47sdRh1dUOLdyC6Q=
github.com/antihax/optional v1.0.0/go.mod h1:uupD/76wgC+ih3iEmQUL+0Ugr19nfwCT1kdvxnR2qWY=
github.com/armon/circbuf v0.0.0-20150827004946-bbbad097214e/go.mod h1:3U/XgcO3hCbHZ8TKRvWD2dDTCfh9M9ya+I9JpbB7O8o=
github.com/armon/consul-api v0.0.0-20180202201655-eb2c6b5be1b6/go.mod h1:grANhF5doyWs3UAsr3K4I6qtAmlQcZDesFNEHPZAzj8=
github.com/armon/go-metrics v0.0.0-20180917152333-f0300d1749da/go.mod h1:Q73ZrmVTwzkszR9V5SSuryQ31EELlFMUz1kKyl939pY=
github.com/armon/go-radix v0.0.0-20180808171621-7fddfc383310/go.mod h1:ufUuZ+zHj4x4TnLV4JWEpy2hxWSpsRywHrMgIH9cCH8=
github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a/go.mod h1:lB+ZfQJz7igIIfQNfa7Ml4HSf2uFQQRzpGGRXenZAgY=
github.com/beorn7/perks v0.0.0-20180321164747-3a771d992973/go.mod h1:Dwedo/Wpr24TaqPxmxbtue+5NUziq4I4S80YR8gNf3Q=
github.com/beorn7/perks v1.0.0/go.mod h1:KWe93zE9D1o94FZ5RNwFwVgaQK1VOXiVxmqh+CedLV8=
github.com/beorn7/perks v1.0.1 h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=
github.com/beorn7/perks v1.0.1/go.mod h1:G2ZrVWU2WbWT9wwq4/hrbKbnv/1ERSJQ0ibhJ6rlkpw=
github.com/bgentry/speakeasy v0.1.0/go.mod h1:+zsyZBPWlz7T6j88CTgSN5bM796AkVf0kBD4zp0CCIs=
github.com/bketelsen/crypt v0.0.4/go.mod h1:aI6NrJ0pMGgvZKL1iVgXLnfIFJtfV+bKCoqOes/6LfM=
github.com/census-instrumentation/opencensus-proto v0.2.1/go.mod h1:f6KPmirojxKA12rnyqOA5BBL4O983OfeGPqjHWSTneU=
github.com/cespare/xxhash v1.1.0 h1:a6HrQnmkObjyL+Gs60czilIUGqrzKutQD6XZog3p+ko=
github.com/cespare/xxhash v1.1.0/go.mod h1:XrSqR1VqqWfGrhpAt58auRo0WTKS1nRRg3ghfAqPWnc=
github.com/cespare/xxhash/v2 v2.1.1 h1:6MnRN8NT7+YBpUIWxHtefFZOKTAPgGjpQSxqLNn0+qY=
github.com/cespare/xxhash/v2 v2.1.1/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/chzyer/logex v1.1.10/go.mod h1:+Ywpsq7O8HXn0nuIou7OrIPyXbp3wmkHB+jjWRnGsAI=
github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e/go.mod h1:nSuG5e5PlCu98SY8svDHJxuZscDgtXS6KTTbou5AhLI=
github.com/chzyer/test v0.0.0-20180213035817-a1ea475d72b1/go.mod h1:Q3SI9o4m/ZMnBNeIyt5eFwwo7qiLfzFZmjNmxjkiQlU=
github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd h1:qMd81Ts1T2OTKmB4acZcyKaMtRnY5Y44NuXGX2GFJ1w=
github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd/go.mod h1:sE/e/2PUdi/liOCUjSTXgM1o87ZssimdTWN964YiIeI=
github.com/coreos/bbolt v1.3.2/go.mod h1:iRUV2dpdMOn7Bo10OQBFzIJO9kkE559Wcmn+qkEiiKk=
github.com/coreos/bbolt v1.3.3 h1:n6AiVyVRKQFNb6mJlwESEvvLoDyiTzXX7ORAUlkeBdY=
github.com/coreos/bbolt v1.3.3/go.mod h1:iRUV2dpdMOn7Bo10OQBFzIJO9kkE559Wcmn+qkEiiKk=
github.com/coreos/etcd v3.3.10+incompatible/go.mod h1:uF7uidLiAD3TWHmW31ZFd/JWoc32PjwdhPthX9715RE=
github.com/coreos/etcd v3.3.18+incompatible h1:Zz1aXgDrFFi1nadh58tA9ktt06cmPTwNNP3dXwIq1lE=
github.com/coreos/etcd v3.3.18+incompatible/go.mod h1:uF7uidLiAD3TWHmW31ZFd/JWoc32PjwdhPthX9715RE=
github.com/coreos/go-oidc/v3 v3.0.0 h1:/mAA0XMgYJw2Uqm7WKGCsKnjitE/+A0FFbOmiRJm7LQ=
github.com/coreos/go-oidc/v3 v3.0.0/go.mod h1:rEJ/idjfUyfkBit1eI1fvyr+64/g9dcKpAm8MJMesvo=
github.com/coreos/go-semver v0.2.0/go.mod h1:nnelYz7RCh+5ahJtPPxZlU+153eP4D4r3EedlOD2RNk=
github.com/coreos/go-semver v0.3.0 h1:wkHLiw0WNATZnSG7epLsujiMCgPAc9xhjJ4tgnAxmfM=
github.com/coreos/go-semver v0.3.0/go.mod h1:nnelYz7RCh+5ahJtPPxZlU+153eP4D4r3EedlOD2RNk=
github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e/go.mod h1:F5haX7vjVVG0kc13fIWeqUViNPyEJxv/OmvnBo0Yme4=
github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf h1:iW4rZ826su+pqaw19uhpSCzhj44qo35pNgKFGqzDKkU=
github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf/go.mod h1:F5haX7vjVVG0kc13fIWeqUViNPyEJxv/OmvnBo0Yme4=
github.com/coreos/go-systemd/v22 v22.3.2/go.mod h1:Y58oyj3AT4RCenI/lSvhwexgC+NSVTIJ3seZv2GcEnc=
github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f h1:lBNOc5arjvs8E5mO2tbpBpLoyyu8B6e44T7hJy6potg=
github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f/go.mod h1:E3G3o1h8I7cfcXa63jLwjI0eiQQMgzzUDFVpN/nH/eA=
github.com/creack/pty v1.1.9/go.mod h1:oKZEueFk5CKHvIhNR5MUki03XCEU+Q6VDXinZuGJ33E=
github.com/davecgh/go-spew v1.1.0/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/dgrijalva/jwt-go v3.2.0+incompatible h1:7qlOGliEKZXTDg6OTjfoBKDXWrumCAMpl/TFQ4/5kLM=
github.com/dgrijalva/jwt-go v3.2.0+incompatible/go.mod h1:E3ru+11k8xSBh+hMPgOLZmtrrCbhqsmaPHjLKYnJCaQ=
github.com/dgryski/go-sip13 v0.0.0-20181026042036-e10d5fee7954/go.mod h1:vAd38F8PWV+bWy6jNmig1y/TA+kYO4g3RSRF0IAv0no=
github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815/go.mod h1:WwZ+bS3ebgob9U8Nd0kOddGdZWjyMGR8Wziv+TBNwSE=
github.com/elazarl/goproxy v0.0.0-20180725130230-947c36da3153/go.mod h1:/Zj4wYkgs4iZTTu3o/KG3Itv/qCCa8VVMlb3i9OVuzc=
github.com/emicklei/go-restful v0.0.0-20170410110728-ff4f55a20633/go.mod h1:otzb+WCGbkyDHkqmQmT5YD2WR4BBwUdeQoFo8l/7tVs=
github.com/envoyproxy/go-control-plane v0.9.1-0.20191026205805-5f8ba28d4473/go.mod h1:YTl/9mNaCwkRvm6d1a2C3ymFceY/DCBVvsKhRF0iEA4=
github.com/envoyproxy/protoc-gen-validate v0.1.0/go.mod h1:iSmxcyjqTsJpI2R4NaDN7+kN2VEUnK/pcBlmesArF7c=
github.com/evanphx/json-patch v4.9.0+incompatible/go.mod h1:50XU6AFN0ol/bzJsmQLiYLvXMP4fmwYFNcr97nuDLSk=
github.com/fatih/color v1.7.0/go.mod h1:Zm6kSWBoL9eyXnKyktHP6abPY2pDugNf5KwzbycvMj4=
github.com/fsnotify/fsnotify v1.4.7/go.mod h1:jwhsz4b93w/PPRr/qN1Yymfu8t87LnFCMoQvtojpjFo=
github.com/fsnotify/fsnotify v1.4.9 h1:hsms1Qyu0jgnwNXIxa+/V/PDsU6CfLf6CNO8H7IWoS4=
github.com/fsnotify/fsnotify v1.4.9/go.mod h1:znqG4EE+3YCdAaPaxE2ZRY/06pZUdp0tY4IgpuI1SZQ=
github.com/ghodss/yaml v1.0.0/go.mod h1:4dBDuWmgqj2HViK6kFavaiC9ZROes6MMH2rRYeMEF04=
github.com/go-gl/glfw v0.0.0-20190409004039-e6da0acd62b1/go.mod h1:vR7hzQXu2zJy9AVAgeJqvqgH9Q5CA+iKCZ2gyEVpxRU=
github.com/go-gl/glfw/v3.3/glfw v0.0.0-20191125211704-12ad95a8df72/go.mod h1:tQ2UAYgL5IevRw8kRxooKSPJfGvJ9fJQFa0TUsXzTg8=
github.com/go-gl/glfw/v3.3/glfw v0.0.0-20200222043503-6f7a984d4dc4/go.mod h1:tQ2UAYgL5IevRw8kRxooKSPJfGvJ9fJQFa0TUsXzTg8=
github.com/go-kit/kit v0.8.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
github.com/go-kit/kit v0.9.0/go.mod h1:xBxKIO96dXMWWy0MnWVtmwkA9/13aqxPnvrjFYMA2as=
github.com/go-logfmt/logfmt v0.3.0/go.mod h1:Qt1PoO58o5twSAckw1HlFXLmHsOX5/0LbT9GBnD5lWE=
github.com/go-logfmt/logfmt v0.4.0/go.mod h1:3RMwSq7FuexP4Kalkev3ejPJsZTpXXBr9+V4qmtdjCk=
github.com/go-logr/logr v0.1.0/go.mod h1:ixOQHD9gLJUVQQ2ZOR7zLEifBX6tGkNJF4QyIY7sIas=
github.com/go-logr/logr v0.4.0 h1:K7/B1jt6fIBQVd4Owv2MqGQClcgf0R266+7C/QjRcLc=
github.com/go-logr/logr v0.4.0/go.mod h1:z6/tIYblkpsD+a4lm/fGIIU9mZ+XfAiaFtq7xTgseGU=
github.com/go-openapi/jsonpointer v0.19.2/go.mod h1:3akKfEdA7DF1sugOqz1dVQHBcuDBPKZGEoHC/NkiQRg=
github.com/go-openapi/jsonpointer v0.19.3/go.mod h1:Pl9vOtqEWErmShwVjC8pYs9cog34VGT37dQOVbmoatg=
github.com/go-openapi/jsonreference v0.19.2/go.mod h1:jMjeRr2HHw6nAVajTXJ4eiUwohSTlpa0o73RUL1owJc=
github.com/go-openapi/jsonreference v0.19.3/go.mod h1:rjx6GuL8TTa9VaixXglHmQmIL98+wF9xc8zWvFonSJ8=
github.com/go-openapi/spec v0.19.3/go.mod h1:FpwSN1ksY1eteniUU7X0N/BgJ7a4WvBFVA8Lj9mJglo=
github.com/go-openapi/swag v0.19.2/go.mod h1:POnQmlKehdgb5mhVOsnJFsivZCEZ/vjK9gh66Z9tfKk=
github.com/go-openapi/swag v0.19.5/go.mod h1:POnQmlKehdgb5mhVOsnJFsivZCEZ/vjK9gh66Z9tfKk=
github.com/go-stack/stack v1.8.0/go.mod h1:v0f6uXyyMGvRgIKkXu+yp6POWl0qKG85gN/melR3HDY=
github.com/godbus/dbus/v5 v5.0.4/go.mod h1:xhWf0FNVPg57R7Z0UbKHbJfkEywrmjJnf7w5xrFpKfA=
github.com/gogo/gateway v1.1.0 h1:u0SuhL9+Il+UbjM9VIE3ntfRujKbvVpFvNB4HbjeVQ0=
github.com/gogo/gateway v1.1.0/go.mod h1:S7rR8FRQyG3QFESeSv4l2WnsyzlCLG0CzBbUUo/mbic=
github.com/gogo/protobuf v1.1.1/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
github.com/gogo/protobuf v1.2.0/go.mod h1:r8qH/GZQm5c6nD/R0oafs1akxWv10x8SbQlK7atdtwQ=
github.com/gogo/protobuf v1.2.1/go.mod h1:hp+jE20tsWTFYpLwKvXlhS1hjn+gTNwPg2I6zVXpSg4=
github.com/gogo/protobuf v1.3.1/go.mod h1:SlYgWuQ5SjCEi6WLHjHCa1yvBfUnHcTbrrZtXPKa29o=
github.com/gogo/protobuf v1.3.2 h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=
github.com/gogo/protobuf v1.3.2/go.mod h1:P1XiOD3dCwIKUDQYPy72D8LYyHL2YPYrpS2s69NZV8Q=
github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=
github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b/go.mod h1:SBH7ygxi8pfUlaOkMMuAQtPIUF8ecWP5IEl/CR7VP2Q=
github.com/golang/groupcache v0.0.0-20190129154638-5b532d6fd5ef/go.mod h1:cIg4eruTrX1D+g88fzRXU5OdNfaM+9IcxsU14FzY7Hc=
github.com/golang/groupcache v0.0.0-20190702054246-869f871628b6/go.mod h1:cIg4eruTrX1D+g88fzRXU5OdNfaM+9IcxsU14FzY7Hc=
github.com/golang/groupcache v0.0.0-20191227052852-215e87163ea7/go.mod h1:cIg4eruTrX1D+g88fzRXU5OdNfaM+9IcxsU14FzY7Hc=
github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e h1:1r7pUrabqp18hOBcwBwiTsbnFeTZHV9eER/QT5JVZxY=
github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e/go.mod h1:cIg4eruTrX1D+g88fzRXU5OdNfaM+9IcxsU14FzY7Hc=
github.com/golang/mock v1.1.1/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
github.com/golang/mock v1.2.0/go.mod h1:oTYuIxOrZwtPieC+H1uAHpcLFnEyAGVDL/k47Jfbm0A=
github.com/golang/mock v1.3.1/go.mod h1:sBzyDLLjw3U8JLTeZvSv8jJB+tU5PVekmnlKIyFUx0Y=
github.com/golang/mock v1.4.0/go.mod h1:UOMv5ysSaYNkG+OFQykRIcU/QvvxJf3p21QfJ2Bt3cw=
github.com/golang/mock v1.4.1/go.mod h1:UOMv5ysSaYNkG+OFQykRIcU/QvvxJf3p21QfJ2Bt3cw=
github.com/golang/mock v1.4.3/go.mod h1:UOMv5ysSaYNkG+OFQykRIcU/QvvxJf3p21QfJ2Bt3cw=
github.com/golang/mock v1.4.4/go.mod h1:l3mdAwkq5BuhzHwde/uurv3sEJeZMXNpwsxVWU71h+4=
github.com/golang/mock v1.5.0/go.mod h1:CWnOUgYIOo4TcNZ0wHX3YZCqsaM1I1Jvs6v3mP3KVu8=
github.com/golang/protobuf v1.3.5 h1:F768QJ1E9tib+q5Sc8MkdJi1RxLTbRcTf8LJV56aRls=
github.com/golang/protobuf v1.3.5/go.mod h1:6O5/vntMXwX2lRkT1hjjk0nAC1IDOTvTlVgjlRvqsdk=
github.com/google/btree v0.0.0-20180813153112-4030bb1f1f0c/go.mod h1:lNA+9X1NB3Zf8V7Ke586lFgjr2dZNuvo3lPJSGZ5JPQ=
github.com/google/btree v1.0.0 h1:0udJVsspx3VBr5FwtLhQQtuAsVc79tTq0ocGIPAU6qo=
github.com/google/btree v1.0.0/go.mod h1:lNA+9X1NB3Zf8V7Ke586lFgjr2dZNuvo3lPJSGZ5JPQ=
github.com/google/go-cmp v0.2.0/go.mod h1:oXzfMopK8JAjlY9xF4vHSVASa0yLyX7SntLO5aqRK0M=
github.com/google/go-cmp v0.3.0/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
github.com/google/go-cmp v0.3.1/go.mod h1:8QqcDgzrUqlUb/G2PQTWiueGozuR1884gddMywk6iLU=
github.com/google/go-cmp v0.4.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.4.1/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.0/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.1/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.2/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.3/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.4/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/go-cmp v0.5.5 h1:Khx7svrCpmxxtHBq5j2mp/xVjsi8hQMfNLvJFAlrGgU=
github.com/google/go-cmp v0.5.5/go.mod h1:v8dTdLbMG2kIc/vJvl+f65V22dbkXbowE6jgT/gNBxE=
github.com/google/gofuzz v1.0.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/google/gofuzz v1.1.0/go.mod h1:dBl0BpW6vV/+mYPU4Po3pmUjxk6FQPldtuIdl/M65Eg=
github.com/google/martian v2.1.0+incompatible/go.mod h1:9I4somxYTbIHy5NJKHRl3wXiIaQGbYVAs8BPL6v8lEs=
github.com/google/martian/v3 v3.0.0/go.mod h1:y5Zk1BBys9G+gd6Jrk0W3cC1+ELVxBWuIGO+w/tUAp0=
github.com/google/martian/v3 v3.1.0/go.mod h1:y5Zk1BBys9G+gd6Jrk0W3cC1+ELVxBWuIGO+w/tUAp0=
github.com/google/pprof v0.0.0-20181206194817-3ea8567a2e57/go.mod h1:zfwlbNMJ+OItoe0UupaVj+oy1omPYYDuagoSzA8v9mc=
github.com/google/pprof v0.0.0-20190515194954-54271f7e092f/go.mod h1:zfwlbNMJ+OItoe0UupaVj+oy1omPYYDuagoSzA8v9mc=
github.com/google/pprof v0.0.0-20191218002539-d4f498aebedc/go.mod h1:ZgVRPoUq/hfqzAqh7sHMqb3I9Rq5C59dIz2SbBwJ4eM=
github.com/google/pprof v0.0.0-20200212024743-f11f1df84d12/go.mod h1:ZgVRPoUq/hfqzAqh7sHMqb3I9Rq5C59dIz2SbBwJ4eM=
github.com/google/pprof v0.0.0-20200229191704-1ebb73c60ed3/go.mod h1:ZgVRPoUq/hfqzAqh7sHMqb3I9Rq5C59dIz2SbBwJ4eM=
github.com/google/pprof v0.0.0-20200430221834-fc25d7d30c6d/go.mod h1:ZgVRPoUq/hfqzAqh7sHMqb3I9Rq5C59dIz2SbBwJ4eM=
github.com/google/pprof v0.0.0-20200708004538-1a94d8640e99/go.mod h1:ZgVRPoUq/hfqzAqh7sHMqb3I9Rq5C59dIz2SbBwJ4eM=
github.com/google/pprof v0.0.0-20201023163331-3e6fc7fc9c4c/go.mod h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=
github.com/google/pprof v0.0.0-20201203190320-1bf35d6f28c2/go.mod h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=
github.com/google/pprof v0.0.0-20210122040257-d980be63207e/go.mod h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=
github.com/google/pprof v0.0.0-20210226084205-cbba55b83ad5/go.mod h1:kpwsk12EmLew5upagYY7GY0pfYCcupk39gWOCRROcvE=
github.com/google/renameio v0.1.0/go.mod h1:KWCgfxg9yswjAJkECMjeO8J8rahYeXnNhOm40UhjYkI=
github.com/google/uuid v1.1.1/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/google/uuid v1.1.2 h1:EVhdT+1Kseyi1/pUmXKaFxYsDNy9RQYkMWRH68J/W7Y=
github.com/google/uuid v1.1.2/go.mod h1:TIyPZe4MgqvfeYDBFedMoGGpEw/LqOeaOT+nhxU+yHo=
github.com/googleapis/gax-go/v2 v2.0.4/go.mod h1:0Wqv26UfaUD9n4G6kQubkQ+KchISgw+vpHVxEJEs9eg=
github.com/googleapis/gax-go/v2 v2.0.5/go.mod h1:DWXyrwAJ9X0FpwwEdw+IPEYBICEFu5mhpdKc/us6bOk=
github.com/googleapis/gnostic v0.4.1/go.mod h1:LRhVm6pbyptWbWbuZ38d1eyptfvIytN3ir6b65WBswg=
github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1 h1:EGx4pi6eqNxGaHF6qqu48+N2wcFQ5qg5FXgOdqsJ5d8=
github.com/gopherjs/gopherjs v0.0.0-20181017120253-0766667cb4d1/go.mod h1:wJfORRmW1u3UXTncJ5qlYoELFm8eSnnEO6hX4iZ3EWY=
github.com/gorilla/websocket v1.4.0/go.mod h1:E7qHFY5m1UJ88s3WnNqhKjPHQ0heANvMoAMk2YaljkQ=
github.com/gorilla/websocket v1.4.1/go.mod h1:YR8l580nyteQvAITg2hZ9XVh4b55+EU/adAjf1fMHhE=
github.com/gorilla/websocket v1.4.2 h1:+/TMaTYc4QFitKJxsQ7Yye35DkWvkdLcvGKqM+x0Ufc=
github.com/gorilla/websocket v1.4.2/go.mod h1:YR8l580nyteQvAITg2hZ9XVh4b55+EU/adAjf1fMHhE=
github.com/grpc-ecosystem/go-grpc-middleware v1.0.0/go.mod h1:FiyG127CGDf3tlThmgyCl78X/SZQqEOJBCDaAfeWzPs=
github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 h1:0IKlLyQ3Hs9nDaiK5cSHAGmcQEIC8l2Ts1u6x5Dfrqg=
github.com/grpc-ecosystem/go-grpc-middleware v1.2.0/go.mod h1:mJzapYve32yjrKlk9GbyCZHuPgZsrbyIbyKhSzOpg6s=
github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 h1:Ovs26xHkKqVztRpIrF/92BcuyuQ/YW4NSIpoGtfXNho=
github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0/go.mod h1:8NvIoxWQoOIhqOTXgfV/d3M/q6VIi02HzZEHgUlZvzk=
github.com/grpc-ecosystem/grpc-gateway v1.8.5/go.mod h1:vNeuVxBJEsws4ogUvrchl83t/GYV9WGTSLVdBhOQFDY=
github.com/grpc-ecosystem/grpc-gateway v1.9.0/go.mod h1:vNeuVxBJEsws4ogUvrchl83t/GYV9WGTSLVdBhOQFDY=
github.com/grpc-ecosystem/grpc-gateway v1.13.0/go.mod h1:8XEsbTttt/W+VvjtQhLACqCisSPWTxCZ7sBRjU6iH9c=
github.com/grpc-ecosystem/grpc-gateway v1.16.0 h1:gmcG1KaJ57LophUzW0Hy8NmPhnMZb4M0+kPpLofRdBo=
github.com/grpc-ecosystem/grpc-gateway v1.16.0/go.mod h1:BDjrQk3hbvj6Nolgz8mAMFbcEtjT1g+wF4CSlocrBnw=
github.com/grpc-kit/api v0.0.0-20200430142535-fa4f5e68cf77 h1:mz4DUdPAX1MnILVxMJQvuzyczgHuk/Op/SMGXpueWJQ=
github.com/grpc-kit/api v0.0.0-20200430142535-fa4f5e68cf77/go.mod h1:Mynyj22C7yDvUVUiRGvKx9R0fI4+AwGTWpZzQgH88A4=
github.com/grpc-kit/cfg v0.1.6 h1:abt3Lo4efMEh8CwtQvAR0mb3/3ZBZ5U8mkktO9qQWPU=
github.com/grpc-kit/cfg v0.1.6/go.mod h1:8xFD3/X/fK2z5LxIXJ8mjG9Ujlz9/mlq7gQ6ZwY0Fhs=
github.com/grpc-kit/pkg v0.0.0-20200515020120-c853412d6e86/go.mod h1:Iptf3Rq5zrbyCfvYMT7KOAurH1M+XyIubGtysKl2DHg=
github.com/grpc-kit/pkg v0.1.1 h1:wzUnrRty4OwXYA4ufHUj1WJiIIz8MNYGIOw0DRuqG7s=
github.com/grpc-kit/pkg v0.1.1/go.mod h1:CFD2zbvkQN+KgR09PJ/eOIQ8E1XWINkWfkFg83RKca8=
github.com/hashicorp/consul/api v1.1.0/go.mod h1:VmuI/Lkw1nC05EYQWNKwWGbkg+FbDBtguAZLlVdkD9Q=
github.com/hashicorp/consul/sdk v0.1.1/go.mod h1:VKf9jXwCTEY1QZP2MOLRhb5i/I/ssyNV1vwHyQBF0x8=
github.com/hashicorp/errwrap v1.0.0/go.mod h1:YH+1FKiLXxHSkmPseP+kNlulaMuP3n2brvKWEqk/Jc4=
github.com/hashicorp/go-cleanhttp v0.5.1/go.mod h1:JpRdi6/HCYpAwUzNwuwqhbovhLtngrth3wmdIIUrZ80=
github.com/hashicorp/go-immutable-radix v1.0.0/go.mod h1:0y9vanUI8NX6FsYoO3zeMjhV/C5i9g4Q3DwcSNZ4P60=
github.com/hashicorp/go-msgpack v0.5.3/go.mod h1:ahLV/dePpqEmjfWmKiqvPkv/twdG7iPBM1vqhUKIvfM=
github.com/hashicorp/go-multierror v1.0.0/go.mod h1:dHtQlpGsu+cZNNAkkCN/P3hoUDHhCYQXV3UM06sGGrk=
github.com/hashicorp/go-rootcerts v1.0.0/go.mod h1:K6zTfqpRlCUIjkwsN4Z+hiSfzSTQa6eBIzfwKfwNnHU=
github.com/hashicorp/go-sockaddr v1.0.0/go.mod h1:7Xibr9yA9JjQq1JpNB2Vw7kxv8xerXegt+ozgdvDeDU=
github.com/hashicorp/go-syslog v1.0.0/go.mod h1:qPfqrKkXGihmCqbJM2mZgkZGvKG1dFdvsLplgctolz4=
github.com/hashicorp/go-uuid v1.0.0/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
github.com/hashicorp/go-uuid v1.0.1/go.mod h1:6SBZvOh/SIDV7/2o3Jml5SYk/TvGqwFJ/bN7x4byOro=
github.com/hashicorp/go.net v0.0.1/go.mod h1:hjKkEWcCURg++eb33jQU7oqQcI9XDCnUzHA0oac0k90=
github.com/hashicorp/golang-lru v0.5.0/go.mod h1:/m3WP610KZHVQ1SGc6re/UDhFvYD7pJ4Ao+sR/qLZy8=
github.com/hashicorp/golang-lru v0.5.1/go.mod h1:/m3WP610KZHVQ1SGc6re/UDhFvYD7pJ4Ao+sR/qLZy8=
github.com/hashicorp/hcl v1.0.0 h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=
github.com/hashicorp/hcl v1.0.0/go.mod h1:E5yfLk+7swimpb2L/Alb/PJmXilQ/rhwaUYs4T20WEQ=
github.com/hashicorp/logutils v1.0.0/go.mod h1:QIAnNjmIWmVIIkWDTG1z5v++HQmx9WQRO+LraFDTW64=
github.com/hashicorp/mdns v1.0.0/go.mod h1:tL+uN++7HEJ6SQLQ2/p+z2pH24WQKWjBPkE0mNTz8vQ=
github.com/hashicorp/memberlist v0.1.3/go.mod h1:ajVTdAv/9Im8oMAAj5G31PhhMCZJV2pPBoIllUwCN7I=
github.com/hashicorp/serf v0.8.2/go.mod h1:6hOLApaqBFA1NXqRQAsxw9QxuDEvNxSQRwA/JwenrHc=
github.com/hpcloud/tail v1.0.0/go.mod h1:ab1qPbhIpdTxEkNHXyeSf5vhxWSCs/tWer42PpOxQnU=
github.com/ianlancetaylor/demangle v0.0.0-20181102032728-5e5cf60278f6/go.mod h1:aSSvb/t6k1mPoxDqO4vJh6VOCGPwU4O0C2/Eqndh1Sc=
github.com/ianlancetaylor/demangle v0.0.0-20200824232613-28f6c0f3b639/go.mod h1:aSSvb/t6k1mPoxDqO4vJh6VOCGPwU4O0C2/Eqndh1Sc=
github.com/jonboulle/clockwork v0.1.0 h1:VKV+ZcuP6l3yW9doeqz6ziZGgcynBVQO+obU0+0hcPo=
github.com/jonboulle/clockwork v0.1.0/go.mod h1:Ii8DK3G1RaLaWxj9trq07+26W01tbo22gdxWY5EU2bo=
github.com/json-iterator/go v1.1.6/go.mod h1:+SdeFBvtyEkXs7REEP0seUULqWtbJapLOCVDaaPEHmU=
github.com/json-iterator/go v1.1.9/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
github.com/json-iterator/go v1.1.10/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
github.com/json-iterator/go v1.1.11 h1:uVUAXhF2To8cbw/3xN3pxj6kk7TYKs98NIrTqPlMWAQ=
github.com/json-iterator/go v1.1.11/go.mod h1:KdQUCv79m/52Kvf8AW2vK1V8akMuk1QjK/uOdHXbAo4=
github.com/jstemmer/go-junit-report v0.0.0-20190106144839-af01ea7f8024/go.mod h1:6v2b51hI/fHJwM22ozAgKL4VKDeJcHhJFhtBdhmNjmU=
github.com/jstemmer/go-junit-report v0.9.1/go.mod h1:Brl9GWCQeLvo8nXZwPNNblvFj/XSXhF0NWZEnDohbsk=
github.com/jtolds/gls v4.20.0+incompatible h1:xdiiI2gbIgH/gLH7ADydsJ1uDOEzR8yvV7C0MuV77Wo=
github.com/jtolds/gls v4.20.0+incompatible/go.mod h1:QJZ7F/aHp+rZTRtaJ1ow/lLfFfVYBRgL+9YlvaHOwJU=
github.com/julienschmidt/httprouter v1.2.0/go.mod h1:SYymIcj16QtmaHHD7aYtjjsJG7VTCxuUUipMqKk8s4w=
github.com/kisielk/errcheck v1.1.0/go.mod h1:EZBBE59ingxPouuu3KfxchcWSUPOHkagtvWXihfKN4Q=
github.com/kisielk/errcheck v1.2.0/go.mod h1:/BMXB+zMLi60iA8Vv6Ksmxu/1UDYcXs4uQLJ+jE2L00=
github.com/kisielk/errcheck v1.5.0/go.mod h1:pFxgyoBC7bSaBwPgfKdkLd5X25qrDl4LWUI2bnpBCr8=
github.com/kisielk/gotool v1.0.0/go.mod h1:XhKaO+MFFWcvkIS/tQcRk01m1F5IRFswLeQ+oQHNcck=
github.com/konsorten/go-windows-terminal-sequences v1.0.1/go.mod h1:T0+1ngSBFLxvqU3pZ+m/2kptfBszLMUkC4ZK/EgS/cQ=
github.com/kr/fs v0.1.0/go.mod h1:FFnZGqtBN9Gxj7eW1uZ42v5BccTP0vu6NEaFoC2HwRg=
github.com/kr/logfmt v0.0.0-20140226030751-b84e30acd515/go.mod h1:+0opPa2QZZtGFBFZlji/RkVcI2GknAs/DXo4wKdlNEc=
github.com/kr/pretty v0.1.0/go.mod h1:dAy3ld7l9f0ibDNOQOHHMYYIIbhfbHSm3C4ZsoJORNo=
github.com/kr/pretty v0.2.0/go.mod h1:ipq/a2n7PKx3OHsz4KJII5eveXtPO4qwEXGdVfWzfnI=
github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
github.com/kr/pty v1.1.5/go.mod h1:9r2w37qlBe7rQ6e1fg1S/9xpWHSnaqNdHD3WcMdbPDA=
github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
github.com/magiconair/properties v1.8.1/go.mod h1:PppfXfuXeibc/6YijjN8zIbojt8czPbwD3XqdrwzmxQ=
github.com/magiconair/properties v1.8.5 h1:b6kJs+EmPFMYGkow9GiUyCyOvIwYetYJ3fSaWak/Gls=
github.com/magiconair/properties v1.8.5/go.mod h1:y3VJvCyxH9uVvJTWEGAELF3aiYNyPKd5NZ3oSwXrF60=
github.com/mailru/easyjson v0.0.0-20190614124828-94de47d64c63/go.mod h1:C1wdFJiN94OJF2b5HbByQZoLdCWB1Yqtg26g4irojpc=
github.com/mailru/easyjson v0.0.0-20190626092158-b2ccc519800e/go.mod h1:C1wdFJiN94OJF2b5HbByQZoLdCWB1Yqtg26g4irojpc=
github.com/mattn/go-colorable v0.0.9/go.mod h1:9vuHe8Xs5qXnSaW/c/ABM9alt+Vo+STaOChaDxuIBZU=
github.com/mattn/go-isatty v0.0.3/go.mod h1:M+lRXTBqGeGNdLjl/ufCoiOlB5xdOkqRJdNxMWT7Zi4=
github.com/matttproud/golang_protobuf_extensions v1.0.1 h1:4hp9jkHxhMHkqkrB3Ix0jegS5sx/RkqARlsWZ6pIwiU=
github.com/matttproud/golang_protobuf_extensions v1.0.1/go.mod h1:D8He9yQNgCq6Z5Ld7szi9bcBfOoFv/3dc6xSMkL2PC0=
github.com/miekg/dns v1.0.14/go.mod h1:W1PPwlIAgtquWBMBEV9nkV9Cazfe8ScdGz/Lj7v3Nrg=
github.com/mitchellh/cli v1.0.0/go.mod h1:hNIlj7HEI86fIcpObd7a0FcrxTWetlwJDGcceTlRvqc=
github.com/mitchellh/go-homedir v1.0.0/go.mod h1:SfyaCUpYCn1Vlf4IUYiD9fPX4A5wJrkLzIz1N1q0pr0=
github.com/mitchellh/go-testing-interface v1.0.0/go.mod h1:kRemZodwjscx+RGhAo8eIhFbs2+BFgRtFPeD/KE+zxI=
github.com/mitchellh/gox v0.4.0/go.mod h1:Sd9lOJ0+aimLBi73mGofS1ycjY8lL3uZM3JPS42BGNg=
github.com/mitchellh/iochan v1.0.0/go.mod h1:JwYml1nuB7xOzsp52dPpHFffvOCDupsG0QubkSMEySY=
github.com/mitchellh/mapstructure v0.0.0-20160808181253-ca63d7c062ee/go.mod h1:FVVH3fgwuzCH5S8UJGiWEs2h04kUh9fWfEaFds41c1Y=
github.com/mitchellh/mapstructure v1.1.2/go.mod h1:FVVH3fgwuzCH5S8UJGiWEs2h04kUh9fWfEaFds41c1Y=
github.com/mitchellh/mapstructure v1.4.1 h1:CpVNEelQCZBooIPDn+AR3NpivK/TIKU8bDxdASFVQag=
github.com/mitchellh/mapstructure v1.4.1/go.mod h1:bFUtVrKA4DC2yAKiSyO/QUcy7e+RRV2QTWOzhPopBRo=
github.com/moby/spdystream v0.2.0/go.mod h1:f7i0iNDQJ059oMTcWxx8MA/zKFIuD/lY+0GqbN2Wy8c=
github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=
github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd/go.mod h1:6dJC0mAP4ikYIbvyc7fijjWJddQyLn8Ig3JB5CqoB9Q=
github.com/modern-go/reflect2 v0.0.0-20180701023420-4b7aa43c6742/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
github.com/modern-go/reflect2 v1.0.1 h1:9f412s+6RmYXLWZSEzVVgPGK7C2PphHj5RJrvfx9AWI=
github.com/modern-go/reflect2 v1.0.1/go.mod h1:bx2lNnkwVCuqBIxFjflWJWanXIb3RllmbCylyMrvgv0=
github.com/munnerz/goautoneg v0.0.0-20120707110453-a547fc61f48d/go.mod h1:+n7T8mK8HuQTcFwEeznm/DIxMOiR9yIdICNftLE1DvQ=
github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223/go.mod h1:qRWi+5nqEBWmkhHvq77mSJWrCKwh8bxhgT7d/eI7P4U=
github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f/go.mod h1:ZdcZmHo+o7JKHSa8/e818NopupXU1YMK5fe1lsApnBw=
github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e h1:fD57ERR4JtEqsWbfPhv4DMiApHyliiK5xCTNVSPiaAs=
github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e/go.mod h1:zD1mROLANZcx1PVRCS0qkT7pwLkGfwJo4zjcN/Tysno=
github.com/oklog/ulid v1.3.1/go.mod h1:CirwcVhetQ6Lv90oh/F+FBtV6XMibvdAFo93nm5qn4U=
github.com/onsi/ginkgo v0.0.0-20170829012221-11459a886d9c/go.mod h1:lLunBs/Ym6LB5Z9jYTR76FiuTmxDTDusOGeTQH+WWjE=
github.com/onsi/ginkgo v1.6.0/go.mod h1:lLunBs/Ym6LB5Z9jYTR76FiuTmxDTDusOGeTQH+WWjE=
github.com/onsi/ginkgo v1.11.0/go.mod h1:lLunBs/Ym6LB5Z9jYTR76FiuTmxDTDusOGeTQH+WWjE=
github.com/onsi/gomega v0.0.0-20170829124025-dcabb60a477c/go.mod h1:C1qb7wdrVGGVU+Z6iS04AVkA3Q65CEZX59MT0QO5uiA=
github.com/onsi/gomega v1.7.0/go.mod h1:ex+gbHU/CVuBBDIJjb2X0qEXbFg53c61hWP/1CpauHY=
github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9 h1:QsgXACQhd9QJhEmRumbsMQQvBtmdS0mafoVEBplWXEg=
github.com/opentracing-contrib/go-stdlib v0.0.0-20190519235532-cf7a6c988dc9/go.mod h1:PLldrQSroqzH70Xl+1DQcGnefIbqsKR7UDaiux3zV+w=
github.com/opentracing/opentracing-go v1.1.0 h1:pWlfV3Bxv7k65HYwkikxat0+s3pV4bsqf19k25Ur8rU=
github.com/opentracing/opentracing-go v1.1.0/go.mod h1:UkNAQd3GIcIGf0SeVgPpRdFStlNbqXla1AfSYxPUl2o=
github.com/pascaldekloe/goe v0.0.0-20180627143212-57f6aae5913c/go.mod h1:lzWF7FIEvWOWxwDKqyGYQf6ZUaNfKdP144TG7ZOy1lc=
github.com/pelletier/go-toml v1.2.0/go.mod h1:5z9KED0ma1S8pY6P1sdut58dfprrGBbd/94hg7ilaic=
github.com/pelletier/go-toml v1.9.3 h1:zeC5b1GviRUyKYd6OJPvBU/mcVDVoL1OhT17FCt5dSQ=
github.com/pelletier/go-toml v1.9.3/go.mod h1:u1nR/EPcESfeI/szUZKdtJ0xRNbUoANCkoOuaOx1Y+c=
github.com/pkg/errors v0.8.0/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
github.com/pkg/errors v0.8.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
github.com/pkg/errors v0.9.1 h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=
github.com/pkg/errors v0.9.1/go.mod h1:bwawxfHBFNV+L2hUp1rHADufV3IMtnDRdf1r5NINEl0=
github.com/pkg/sftp v1.10.1/go.mod h1:lYOWFsE0bwd1+KfKJaKeuokY15vzFx25BLbzYYoAxZI=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/posener/complete v1.1.1/go.mod h1:em0nMJCgc9GFtwrmVmEMR/ZL6WyhyjMBndrE9hABlRI=
github.com/prometheus/client_golang v0.9.1/go.mod h1:7SWBe2y4D6OKWSNQJUaRYU/AaXPKyh/dDVn+NZz0KFw=
github.com/prometheus/client_golang v0.9.3/go.mod h1:/TN21ttK/J9q6uSwhBd54HahCDft0ttaMvbicHlPoso=
github.com/prometheus/client_golang v1.0.0/go.mod h1:db9x61etRT2tGnBNRi70OPL5FsnadC4Ky3P0J6CfImo=
github.com/prometheus/client_golang v1.4.1 h1:FFSuS004yOQEtDdTq+TAOLP5xUq63KqAFYyOi8zA+Y8=
github.com/prometheus/client_golang v1.4.1/go.mod h1:e9GMxYsXl05ICDXkRhurwBS4Q3OK1iX/F2sw+iXX5zU=
github.com/prometheus/client_model v0.0.0-20180712105110-5c3871d89910/go.mod h1:MbSGuTsp3dbXC40dX6PRTWyKYBIrTGTE9sqQNg2J8bo=
github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
github.com/prometheus/client_model v0.0.0-20190812154241-14fe0d1b01d4/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
github.com/prometheus/client_model v0.2.0 h1:uq5h0d+GuxiXLJLNABMgp2qUWDPiLvgCzz2dUR+/W/M=
github.com/prometheus/client_model v0.2.0/go.mod h1:xMI15A0UPsDsEKsMN9yxemIoYk6Tm2C1GtYGdfGttqA=
github.com/prometheus/common v0.0.0-20181113130724-41aa239b4cce/go.mod h1:daVV7qP5qjZbuso7PdcryaAu0sAZbrN9i7WWcTMWvro=
github.com/prometheus/common v0.4.0/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
github.com/prometheus/common v0.4.1/go.mod h1:TNfzLD0ON7rHzMJeJkieUDPYmFC7Snx/y86RQel1bk4=
github.com/prometheus/common v0.9.1 h1:KOMtN28tlbam3/7ZKEYKHhKoJZYYj3gMH4uc62x7X7U=
github.com/prometheus/common v0.9.1/go.mod h1:yhUN8i9wzaXS3w1O07YhxHEBxD+W35wd8bs7vj7HSQ4=
github.com/prometheus/procfs v0.0.0-20181005140218-185b4288413d/go.mod h1:c3At6R/oaqEKCNdg8wHV1ftS6bRYblBhIjjI8uT2IGk=
github.com/prometheus/procfs v0.0.0-20190507164030-5867b95ac084/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
github.com/prometheus/procfs v0.0.2/go.mod h1:TjEm7ze935MbeOT/UhFTIMYKhuLP4wbCsTZCD3I8kEA=
github.com/prometheus/procfs v0.0.8 h1:+fpWZdT24pJBiqJdAwYBjPSk+5YmQzYNPYzQsdzLkt8=
github.com/prometheus/procfs v0.0.8/go.mod h1:7Qr8sr6344vo1JqZ6HhLceV9o3AJ1Ff+GxbHq6oeK9A=
github.com/prometheus/tsdb v0.7.1/go.mod h1:qhTCs0VvXwvX/y3TZrWD7rabWM+ijKTux40TwIPHuXU=
github.com/rogpeppe/fastuuid v0.0.0-20150106093220-6724a57986af/go.mod h1:XWv6SoW27p1b0cqNHllgS5HIMJraePCO15w5zCzIWYg=
github.com/rogpeppe/fastuuid v1.2.0/go.mod h1:jVj6XXZzXRy/MSR5jhDC/2q6DgLz+nrA6LYCDYWNEvQ=
github.com/rogpeppe/go-internal v1.3.0/go.mod h1:M8bDsm7K2OlrFYOpmOWEs/qY81heoFRclV5y23lUDJ4=
github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f/go.mod h1:sm1tb6uqfes/u+d4ooFouqFdy9/2g9QGwK3SQygK0Ts=
github.com/sean-/seed v0.0.0-20170313163322-e2103e2c3529/go.mod h1:DxrIzT+xaE7yg65j358z/aeFdxmN0P9QXhEzd20vsDc=
github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 h1:bUGsEnyNbVPw06Bs80sCeARAlK8lhwqGyi6UT8ymuGk=
github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749/go.mod h1:ZY1cvUeJuFPAdZ/B6v7RHavJWZn2YPVFQ1OSXhCGOkg=
github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 h1:pXY9qYc/MP5zdvqWEUH6SjNiu7VhSjuVFTFiTcphaLU=
github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546/go.mod h1:TrYk7fJVaAttu97ZZKrO9UbRa8izdowaMIZcxYMbVaw=
github.com/sirupsen/logrus v1.2.0/go.mod h1:LxeOpSwHxABJmUn/MG1IvRgCAasNZTLOkJPxbbu5VWo=
github.com/sirupsen/logrus v1.4.2/go.mod h1:tLMulIdttU9McNUspp0xgXVQah82FyeX6MwdIuYE2rE=
github.com/sirupsen/logrus v1.8.1 h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=
github.com/sirupsen/logrus v1.8.1/go.mod h1:yWOB1SBYBC5VeMP7gHvWumXLIWorT60ONWic61uBYv0=
github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d h1:zE9ykElWQ6/NYmHa3jpm/yHnI4xSofP+UP6SpjHcSeM=
github.com/smartystreets/assertions v0.0.0-20180927180507-b2de0cb4f26d/go.mod h1:OnSkiWE9lh6wB0YB77sQom3nweQdgAjqCqsofrRNTgc=
github.com/smartystreets/goconvey v1.6.4 h1:fv0U8FUIMPNf1L9lnHLvLhgicrIVChEkdzIKYqbNC9s=
github.com/smartystreets/goconvey v1.6.4/go.mod h1:syvi0/a8iFYH4r/RixwvyeAJjdLS9QV7WQ/tjFTllLA=
github.com/soheilhy/cmux v0.1.4 h1:0HKaf1o97UwFjHH9o5XsHUOF+tqmdA7KEzXLpiyaw0E=
github.com/soheilhy/cmux v0.1.4/go.mod h1:IM3LyeVVIOuxMH7sFAkER9+bJ4dT7Ms6E4xg4kGIyLM=
github.com/spaolacci/murmur3 v0.0.0-20180118202830-f09979ecbc72/go.mod h1:JwIasOWyU6f++ZhiEuf87xNszmSA2myDM2Kzu9HwQUA=
github.com/spf13/afero v1.1.2/go.mod h1:j4pytiNVoe2o6bmDsKpLACNPDBIoEAkihy7loJ1B0CQ=
github.com/spf13/afero v1.6.0 h1:xoax2sJ2DT8S8xA2paPFjDCScCNeWsg75VG0DLRreiY=
github.com/spf13/afero v1.6.0/go.mod h1:Ai8FlHk4v/PARR026UzYexafAt9roJ7LcLMAmO6Z93I=
github.com/spf13/cast v1.3.0/go.mod h1:Qx5cxh0v+4UWYiBimWS+eyWzqEqokIECu5etghLkUJE=
github.com/spf13/cast v1.3.1 h1:nFm6S0SMdyzrzcmThSipiEubIDy8WEXKNZ0UOgiRpng=
github.com/spf13/cast v1.3.1/go.mod h1:Qx5cxh0v+4UWYiBimWS+eyWzqEqokIECu5etghLkUJE=
github.com/spf13/jwalterweatherman v1.0.0/go.mod h1:cQK4TGJAtQXfYWX+Ddv3mKDzgVb68N+wFjFa4jdeBTo=
github.com/spf13/jwalterweatherman v1.1.0 h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=
github.com/spf13/jwalterweatherman v1.1.0/go.mod h1:aNWZUN0dPAAO/Ljvb5BEdw96iTZ0EXowPYD95IqWIGo=
github.com/spf13/pflag v0.0.0-20170130214245-9ff6c6923cff/go.mod h1:DYY7MBk1bdzusC3SYhjObp+wFpr4gzcvqqNjLnInEg4=
github.com/spf13/pflag v1.0.3/go.mod h1:DYY7MBk1bdzusC3SYhjObp+wFpr4gzcvqqNjLnInEg4=
github.com/spf13/pflag v1.0.5 h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=
github.com/spf13/pflag v1.0.5/go.mod h1:McXfInJRrz4CZXVZOBLb0bTZqETkiAhM9Iw0y3An2Bg=
github.com/spf13/viper v1.6.2/go.mod h1:t3iDnF5Jlj76alVNuyFBk5oUMCvsrkbvZK0WQdfDi5k=
github.com/spf13/viper v1.8.1 h1:Kq1fyeebqsBfbjZj4EL7gj2IO0mMaiyjYUWcUsl2O44=
github.com/spf13/viper v1.8.1/go.mod h1:o0Pch8wJ9BVSWGQMbra6iw0oQ5oktSIBaujf1rJH9Ns=
github.com/stretchr/objx v0.1.0/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.1.1/go.mod h1:HFkY916IF+rwdDfMAkV7OtwuqBVzrE8GR6GFx+wExME=
github.com/stretchr/objx v0.2.0/go.mod h1:qt09Ya8vawLte6SNmTgCsAVtYtaKzEcn8ATUoHMkEqE=
github.com/stretchr/testify v1.2.2/go.mod h1:a8OnRcib4nhh0OaRAV+Yts87kKdq0PP7pXfy6kDkUVs=
github.com/stretchr/testify v1.3.0/go.mod h1:M5WIy9Dh21IEIfnGCwXGc5bZfKNJtfHm1UVUgZn+9EI=
github.com/stretchr/testify v1.4.0/go.mod h1:j7eGeouHqKxXV5pUuKE4zz7dFj8WfuZ+81PSLYec5m4=
github.com/stretchr/testify v1.5.1/go.mod h1:5W2xD1RspED5o8YsWQXVCued0rvSQ+mT+I5cxcmMvtA=
github.com/stretchr/testify v1.6.1/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/stretchr/testify v1.7.0 h1:nwc3DEeHmmLAfoZucVR881uASk0Mfjw8xYJ99tb5CcY=
github.com/stretchr/testify v1.7.0/go.mod h1:6Fq8oRcR53rry900zMqJjRRixrwX3KX962/h/Wwjteg=
github.com/subosito/gotenv v1.2.0 h1:Slr1R9HxAlEKefgq5jn9U+DnETlIUa6HfgEzj0g5d7s=
github.com/subosito/gotenv v1.2.0/go.mod h1:N0PQaV/YGNqwC0u51sEeR/aUtSLEXKX9iv69rRypqCw=
github.com/tmc/grpc-websocket-proxy v0.0.0-20190109142713-0ad062ec5ee5/go.mod h1:ncp9v5uamzpCO7NfCPTXjqaC+bZgJeR0sMTm6dMHP7U=
github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc h1:yUaosFVTJwnltaHbSNC3i82I92quFs+OFPRl8kNMVwo=
github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc/go.mod h1:ncp9v5uamzpCO7NfCPTXjqaC+bZgJeR0sMTm6dMHP7U=
github.com/uber/jaeger-client-go v2.22.1+incompatible h1:NHcubEkVbahf9t3p75TOCR83gdUHXjRJvjoBh1yACsM=
github.com/uber/jaeger-client-go v2.22.1+incompatible/go.mod h1:WVhlPFC8FDjOFMMWRy2pZqQJSXxYSwNYOkTr/Z6d3Kk=
github.com/uber/jaeger-lib v2.2.0+incompatible h1:MxZXOiR2JuoANZ3J6DE/U0kSFv/eJ/GfSYVCjK7dyaw=
github.com/uber/jaeger-lib v2.2.0+incompatible/go.mod h1:ComeNDZlWwrWnDv8aPp0Ba6+uUTzImX/AauajbLI56U=
github.com/ugorji/go v1.1.4/go.mod h1:uQMGLiO92mf5W77hV/PUCpI3pbzQx3CRekS0kk+RGrc=
github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 h1:eY9dn8+vbi4tKz5Qo6v2eYzo7kUS51QINcR5jNpbZS8=
github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2/go.mod h1:UETIi67q53MR2AWcXfiuqkDkRtnGDLqkBTpCHuJHxtU=
github.com/xordataexchange/crypt v0.0.3-0.20170626215501-b2862e3d0a77/go.mod h1:aYKd//L2LvnjZzWKhF00oedf4jCCReLcmhLdhm1A27Q=
github.com/yuin/goldmark v1.1.25/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
github.com/yuin/goldmark v1.1.27/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
github.com/yuin/goldmark v1.1.32/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
github.com/yuin/goldmark v1.2.1/go.mod h1:3hX8gzYuyVAZsxl0MRgGTJEmQBFcNTphYh9decYSb74=
go.etcd.io/bbolt v1.3.2/go.mod h1:IbVyRI1SCnLcuJnV2u8VeU0CEYM7e686BmAb1XKL+uU=
go.etcd.io/bbolt v1.3.3 h1:MUGmc65QhB3pIlaQ5bB4LwqSj6GIonVJXpZiaKNyaKk=
go.etcd.io/bbolt v1.3.3/go.mod h1:IbVyRI1SCnLcuJnV2u8VeU0CEYM7e686BmAb1XKL+uU=
go.etcd.io/etcd/api/v3 v3.5.0/go.mod h1:cbVKeC6lCfl7j/8jBhAK6aIYO9XOjdptoxU/nLQcPvs=
go.etcd.io/etcd/client/pkg/v3 v3.5.0/go.mod h1:IJHfcCEKxYu1Os13ZdwCwIUTUVGYTSAM3YSwc9/Ac1g=
go.etcd.io/etcd/client/v2 v2.305.0/go.mod h1:h9puh54ZTgAKtEbut2oe9P4L/oqKCVB6xsXlzd7alYQ=
go.opencensus.io v0.21.0/go.mod h1:mSImk1erAIZhrmZN+AvHh14ztQfjbGwt4TtuofqLduU=
go.opencensus.io v0.22.0/go.mod h1:+kGneAE2xo2IficOXnaByMWTGM9T73dGwxeWcUqIpI8=
go.opencensus.io v0.22.2/go.mod h1:yxeiOL68Rb0Xd1ddK5vPZ/oVn4vY4Ynel7k9FzqtOIw=
go.opencensus.io v0.22.3/go.mod h1:yxeiOL68Rb0Xd1ddK5vPZ/oVn4vY4Ynel7k9FzqtOIw=
go.opencensus.io v0.22.4/go.mod h1:yxeiOL68Rb0Xd1ddK5vPZ/oVn4vY4Ynel7k9FzqtOIw=
go.opencensus.io v0.22.5/go.mod h1:5pWMHQbX5EPX2/62yrJeAkowc+lfs/XD7Uxpq3pI6kk=
go.opencensus.io v0.23.0/go.mod h1:XItmlyltB5F7CS4xOC1DcqMoFqwtC6OG2xF7mCv7P7E=
go.uber.org/atomic v1.4.0/go.mod h1:gD2HeocX3+yG+ygLZcrzQJaqmWj9AIm7n08wl/qW/PE=
go.uber.org/atomic v1.5.0/go.mod h1:sABNBOSYdrvTF6hTgEIbc7YasKWGhgEQZyfxyTvoXHQ=
go.uber.org/atomic v1.7.0 h1:ADUqmZGgLDDfbSL9ZmPxKTybcoEYHgpYfELNoN+7hsw=
go.uber.org/atomic v1.7.0/go.mod h1:fEN4uk6kAWBTFdckzkM89CLk9XfWZrxpCo0nPH17wJc=
go.uber.org/multierr v1.1.0/go.mod h1:wR5kodmAFQ0UK8QlbwjlSNy0Z68gJhDJUG5sjR94q/0=
go.uber.org/multierr v1.3.0/go.mod h1:VgVr7evmIr6uPjLBxg28wmKNXyqE9akIJ5XnfpiKl+4=
go.uber.org/multierr v1.6.0 h1:y6IPFStTAIT5Ytl7/XYmHvzXQ7S3g/IeZW9hyZ5thw4=
go.uber.org/multierr v1.6.0/go.mod h1:cdWPpRnG4AhwMwsgIHip0KRBQjJy5kYEpYjJxpXp9iU=
go.uber.org/tools v0.0.0-20190618225709-2cfd321de3ee/go.mod h1:vJERXedbb3MVM5f9Ejo0C68/HhF8uaILCdgjnY+goOA=
go.uber.org/zap v1.10.0/go.mod h1:vwi/ZaCAaUcBkycHslxD9B2zi4UTXhF60s6SWpuDF0Q=
go.uber.org/zap v1.13.0/go.mod h1:zwrFLgMcdUuIBviXEYEH1YKNaOBnKXsx2IPda5bBwHM=
go.uber.org/zap v1.17.0 h1:MTjgFu6ZLKvY6Pvaqk97GlxNBuMpV4Hy/3P6tRGlI2U=
go.uber.org/zap v1.17.0/go.mod h1:MXVU+bhUf/A7Xi2HNOnopQOrmycQ5Ih87HtOu4q5SSo=
golang.org/x/crypto v0.0.0-20180904163835-0709b304e793/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
golang.org/x/crypto v0.0.0-20181029021203-45a5f77698d3/go.mod h1:6SG95UA2DQfeDnfUPMdvaQW0Q7yPrPDi9nlGo2tz2b4=
golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=
golang.org/x/crypto v0.0.0-20190510104115-cbcb75029529/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
golang.org/x/crypto v0.0.0-20190611184440-5c40567a22f8/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550/go.mod h1:yigFU9vqHzYiE8UmvKecakEJjdnWj3jj499lnFckfCI=
golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 h1:psW17arqaxU48Z5kZ0CQnkZWQJsqcURM6tKiBApRjXI=
golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9/go.mod h1:LzIPMQfyMNhhGPhUkYOs5KpL4U8rLKemX1yGLhDgUto=
golang.org/x/exp v0.0.0-20190121172915-509febef88a4/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
golang.org/x/exp v0.0.0-20190306152737-a1d7652674e8/go.mod h1:CJ0aWSM057203Lf6IL+f9T1iT9GByDxfZKAQTCR3kQA=
golang.org/x/exp v0.0.0-20190510132918-efd6b22b2522/go.mod h1:ZjyILWgesfNpC6sMxTJOJm9Kp84zZh5NQWvqDGG3Qr8=
golang.org/x/exp v0.0.0-20190829153037-c13cbed26979/go.mod h1:86+5VVa7VpoJ4kLfm080zCjGlMRFzhUhsZKEZO7MGek=
golang.org/x/exp v0.0.0-20191030013958-a1ab85dbe136/go.mod h1:JXzH8nQsPlswgeRAPE3MuO9GYsAcnJvJ4vnMwN/5qkY=
golang.org/x/exp v0.0.0-20191129062945-2f5052295587/go.mod h1:2RIsYlXP63K8oxa1u096TMicItID8zy7Y6sNkU49FU4=
golang.org/x/exp v0.0.0-20191227195350-da58074b4299/go.mod h1:2RIsYlXP63K8oxa1u096TMicItID8zy7Y6sNkU49FU4=
golang.org/x/exp v0.0.0-20200119233911-0405dc783f0a/go.mod h1:2RIsYlXP63K8oxa1u096TMicItID8zy7Y6sNkU49FU4=
golang.org/x/exp v0.0.0-20200207192155-f17229e696bd/go.mod h1:J/WKrq2StrnmMY6+EHIKF9dgMWnmCNThgcyBT1FY9mM=
golang.org/x/exp v0.0.0-20200224162631-6cc2880d07d6/go.mod h1:3jZMyOhIsHpP37uCMkUooju7aAi5cS1Q23tOzKc+0MU=
golang.org/x/image v0.0.0-20190227222117-0694c2d4d067/go.mod h1:kZ7UVZpmo3dzQBMxlp+ypCbDeSB+sBbTgSJuh5dn5js=
golang.org/x/image v0.0.0-20190802002840-cff245a6509b/go.mod h1:FeLwcggjj3mMvU+oOTbSwawSJRM1uh48EjtB4UJZlP0=
golang.org/x/lint v0.0.0-20190301231843-5614ed5bae6f/go.mod h1:UVdnD1Gm6xHRNCYTkRU2/jEulfH38KcIWyp/GAMgvoE=
golang.org/x/lint v0.0.0-20190313153728-d0100b6bd8b3/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
golang.org/x/lint v0.0.0-20190409202823-959b441ac422/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
golang.org/x/lint v0.0.0-20190909230951-414d861bb4ac/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
golang.org/x/lint v0.0.0-20190930215403-16217165b5de/go.mod h1:6SW0HCj/g11FgYtHlgUYUwCkIfeOF89ocIRzGO/8vkc=
golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f/go.mod h1:5qLYkcX4OjUUV8bRuDixDT3tpyyb+LUpUlRWLxfhWrs=
golang.org/x/lint v0.0.0-20200130185559-910be7a94367/go.mod h1:3xt1FjdF8hUf6vQPIChWIBhFzV8gjjsPE/fR3IyQdNY=
golang.org/x/lint v0.0.0-20200302205851-738671d3881b/go.mod h1:3xt1FjdF8hUf6vQPIChWIBhFzV8gjjsPE/fR3IyQdNY=
golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5/go.mod h1:3xt1FjdF8hUf6vQPIChWIBhFzV8gjjsPE/fR3IyQdNY=
golang.org/x/mobile v0.0.0-20190312151609-d3739f865fa6/go.mod h1:z+o9i4GpDbdi3rU15maQ/Ox0txvL9dWGYEHz965HBQE=
golang.org/x/mobile v0.0.0-20190719004257-d2bd2a29d028/go.mod h1:E/iHnbuqvinMTCcRqshq8CkpyQDoeVncDDYHnLhea+o=
golang.org/x/mod v0.0.0-20190513183733-4bf6d317e70e/go.mod h1:mXi4GBBbnImb6dmsKGUJ2LatrhH/nqhxcFungHvyanc=
golang.org/x/mod v0.1.0/go.mod h1:0QHyrYULN0/3qlju5TqG8bIK38QM8yzMo5ekMj3DlcY=
golang.org/x/mod v0.1.1-0.20191105210325-c90efee705ee/go.mod h1:QqPTAvyqsEbceGzBzNggFXnrqF1CaUcvgkdR5Ot7KZg=
golang.org/x/mod v0.1.1-0.20191107180719-034126e5016b/go.mod h1:QqPTAvyqsEbceGzBzNggFXnrqF1CaUcvgkdR5Ot7KZg=
golang.org/x/mod v0.2.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
golang.org/x/mod v0.3.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
golang.org/x/mod v0.4.0/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
golang.org/x/mod v0.4.1/go.mod h1:s0Qsj1ACt9ePp/hMypM3fl4fZqREWJwdYDEqhRiZZUA=
golang.org/x/net v0.0.0-20180724234803-3673e40ba225/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
golang.org/x/net v0.0.0-20180906233101-161cd47e91fd/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
golang.org/x/net v0.0.0-20181023162649-9b4f9f5ad519/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
golang.org/x/net v0.0.0-20181114220301-adae6a3d119a/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
golang.org/x/net v0.0.0-20181201002055-351d144fa1fc/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
golang.org/x/net v0.0.0-20181220203305-927f97764cc3/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
golang.org/x/net v0.0.0-20190108225652-1e06a53dbb7e/go.mod h1:mL1N/T3taQHkDXs73rZJwtUhF3w3ftmwwsq0BUmARs4=
golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
golang.org/x/net v0.0.0-20190501004415-9ce7a6920f09/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=
golang.org/x/net v0.0.0-20190522155817-f3200d17e092/go.mod h1:HSz+uSET+XFnRR8LxR5pz3Of3rY3CfYBVs4xY44aLks=
golang.org/x/net v0.0.0-20190603091049-60506f45cf65/go.mod h1:HSz+uSET+XFnRR8LxR5pz3Of3rY3CfYBVs4xY44aLks=
golang.org/x/net v0.0.0-20190613194153-d28f0bde5980/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20190620200207-3b0461eec859/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20190628185345-da137c7871d7/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20190724013045-ca1201d0de80/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20191002035440-2ec189313ef0/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20200202094626-16171245cfb2/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20200222125558-5a598a2470a0/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20200226121028-0de0cce0169b/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20200301022130-244492dfa37a/go.mod h1:z5CRVTTTmAJ677TzLLGU+0bjPO0LkuOLi4/5GtJWs/s=
golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
golang.org/x/net v0.0.0-20200501053045-e0ff5e5a1de5/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
golang.org/x/net v0.0.0-20200505041828-1ed23360d12c/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
golang.org/x/net v0.0.0-20200506145744-7e3656a0809f/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
golang.org/x/net v0.0.0-20200513185701-a91f0712d120/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2/go.mod h1:qpuaurCH72eLCgpAm/N6yyVIVM9cpaDIP3A8BGJEC5A=
golang.org/x/net v0.0.0-20200625001655-4c5254603344/go.mod h1:/O7V0waA8r7cgGh81Ro3o1hOxt32SMVPicZroKQ2sZA=
golang.org/x/net v0.0.0-20200707034311-ab3426394381/go.mod h1:/O7V0waA8r7cgGh81Ro3o1hOxt32SMVPicZroKQ2sZA=
golang.org/x/net v0.0.0-20200822124328-c89045814202/go.mod h1:/O7V0waA8r7cgGh81Ro3o1hOxt32SMVPicZroKQ2sZA=
golang.org/x/net v0.0.0-20201021035429-f5854403a974/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
golang.org/x/net v0.0.0-20201031054903-ff519b6c9102/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
golang.org/x/net v0.0.0-20201110031124-69a78807bb2b/go.mod h1:sp8m0HH+o8qH0wwXwYZr8TS3Oi6o0r6Gce1SSxlDquU=
golang.org/x/net v0.0.0-20201209123823-ac852fbbde11/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
golang.org/x/net v0.0.0-20210119194325-5f4716e94777/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
golang.org/x/net v0.0.0-20210224082022-3d97a244fca7/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
golang.org/x/net v0.0.0-20210226172049-e18ecbb05110/go.mod h1:m0MpNAwzfU5UDzcl9v0D8zg8gWTRqZa9RBIspLL5mdg=
golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4 h1:b0LrWgu8+q7z4J+0Y3Umo5q1dL7NXBkKBWkaVkAq17E=
golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4/go.mod h1:RBQZq4jEuRlivfhVLdyRGr576XBO4/greRjx4P4O3yc=
golang.org/x/oauth2 v0.0.0-20180821212333-d2e6202438be/go.mod h1:N/0e6XlmueqKjAGxoOufVs8QHGRruUQn6yWY3a++T0U=
golang.org/x/oauth2 v0.0.0-20190226205417-e64efc72b421/go.mod h1:gOpvHmFTYa4IltrdGE7lF6nIHvwfUNPOp7c8zoXwtLw=
golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45/go.mod h1:gOpvHmFTYa4IltrdGE7lF6nIHvwfUNPOp7c8zoXwtLw=
golang.org/x/oauth2 v0.0.0-20191202225959-858c2ad4c8b6/go.mod h1:gOpvHmFTYa4IltrdGE7lF6nIHvwfUNPOp7c8zoXwtLw=
golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d/go.mod h1:gOpvHmFTYa4IltrdGE7lF6nIHvwfUNPOp7c8zoXwtLw=
golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43/go.mod h1:KelEdhl1UZF7XfJ4dDtk6s++YSgaE7mD/BuKKDLBl4A=
golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58/go.mod h1:KelEdhl1UZF7XfJ4dDtk6s++YSgaE7mD/BuKKDLBl4A=
golang.org/x/oauth2 v0.0.0-20201208152858-08078c50e5b5/go.mod h1:KelEdhl1UZF7XfJ4dDtk6s++YSgaE7mD/BuKKDLBl4A=
golang.org/x/oauth2 v0.0.0-20210218202405-ba52d332ba99/go.mod h1:KelEdhl1UZF7XfJ4dDtk6s++YSgaE7mD/BuKKDLBl4A=
golang.org/x/oauth2 v0.0.0-20210220000619-9bb904979d93/go.mod h1:KelEdhl1UZF7XfJ4dDtk6s++YSgaE7mD/BuKKDLBl4A=
golang.org/x/oauth2 v0.0.0-20210313182246-cd4f82c27b84/go.mod h1:KelEdhl1UZF7XfJ4dDtk6s++YSgaE7mD/BuKKDLBl4A=
golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602 h1:0Ja1LBD+yisY6RWM/BH7TJVXWsSjs2VwBSmvSX4HdBc=
golang.org/x/oauth2 v0.0.0-20210402161424-2e8d93401602/go.mod h1:KelEdhl1UZF7XfJ4dDtk6s++YSgaE7mD/BuKKDLBl4A=
golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20181108010431-42b317875d0f/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20190227155943-e225da77a7e6/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20190423024810-112230192c58/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20201020160332-67f06af15bc9/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20201207232520-09787c993a3a/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sync v0.0.0-20210220032951-036812b2e83c/go.mod h1:RxMgew5VJxzue5/jJTE5uejpjVlOe/izrB70Jof72aM=
golang.org/x/sys v0.0.0-20180823144017-11551d06cbcc/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20180909124046-d0be0721c37e/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20181026203630-95b1ffbd15a5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20181107165924-66b7b1311ac8/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20181116152217-5ac8a444bdc5/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=
golang.org/x/sys v0.0.0-20190312061237-fead79001313/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190412213103-97732733099d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190422165155-953cdadca894/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190502145724-3ef323f4f1fd/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190507160741-ecd444e8653b/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190606165138-5da285871e9c/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190616124812-15dcb6c0061f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190624142023-c5567b49c5d0/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20190726091711-fc99dfbffb4e/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20191001151750-bb3f8db39f24/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20191005200804-aed5e4c7ecf9/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20191026070338-33540a1f6037/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20191204072324-ce4227a45e2e/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20191228213918-04cbcbbfeed8/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200113162924-86b910548bc1/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200122134326-e047566fdf82/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200202164722-d101bd2416d5/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200212091648-12a6c2dcc1e4/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200223170610-d5e6a3e2c0ae/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200331124033-c3d80250170d/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200501052902-10377860bb8e/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200511232937-7e40ca221e25/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200515095857-1151b9dac4a9/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200523222454-059865788121/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200803210538-64077c9b5642/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200905004654-be1d3432aa8f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20201119102817-f84b799fce68/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20201201145000-ef89a241ccb3/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210104204734-6f8348627aad/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210220050731-9a76102bfb43/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210305230114-8fe3ee5dd75b/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210315160823-c6e025ad8005/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57 h1:F5Gozwx4I1xtr/sr/8CFbb57iKi3297KFs0QDbGN60A=
golang.org/x/sys v0.0.0-20210403161142-5e06dd20ab57/go.mod h1:h1NjWce9XRLGQEsW7wpKNCjG9DtNlClVuFLEZdDNbEs=
golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1/go.mod h1:bj7SfCRtBDWHUb9snDiAeCFNEtKQo2Wmx5Cou7ajbmo=
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.1-0.20180807135948-17ff2d5776d2/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
golang.org/x/text v0.3.2/go.mod h1:bEr9sfX3Q8Zfm5fL9x+3itogRgK3+ptLWKqgva+5dAk=
golang.org/x/text v0.3.3/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/text v0.3.4/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/text v0.3.5 h1:i6eZZ+zk0SOf0xgBpEpPD18qWcJda6q1sxt3S0kzyUQ=
golang.org/x/text v0.3.5/go.mod h1:5Zoc/QRtKVWzQhOtBMvqHzDpF6irO9z98xDceosuGiQ=
golang.org/x/time v0.0.0-20181108054448-85acf8d2951c/go.mod h1:tRJNPiyCQ0inRvYxbN9jk5I+vvW/OXSQhTDSoE431IQ=
golang.org/x/time v0.0.0-20190308202827-9d24e82272b4/go.mod h1:tRJNPiyCQ0inRvYxbN9jk5I+vvW/OXSQhTDSoE431IQ=
golang.org/x/time v0.0.0-20191024005414-555d28b269f0 h1:/5xXl8Y5W96D+TtHSlonuFqGHIWVuyCkGJLwGh9JJFs=
golang.org/x/time v0.0.0-20191024005414-555d28b269f0/go.mod h1:tRJNPiyCQ0inRvYxbN9jk5I+vvW/OXSQhTDSoE431IQ=
golang.org/x/tools v0.0.0-20180221164845-07fd8470d635/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
golang.org/x/tools v0.0.0-20180917221912-90fa682c2a6e/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
golang.org/x/tools v0.0.0-20181030221726-6c7e314b6563/go.mod h1:n7NCudcB/nEzxVGmLbDWY5pfWTLqBcC2KZ6jyYvM4mQ=
golang.org/x/tools v0.0.0-20190311212946-11955173bddd/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
golang.org/x/tools v0.0.0-20190312151545-0bb0c0a6e846/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
golang.org/x/tools v0.0.0-20190312170243-e65039ee4138/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
golang.org/x/tools v0.0.0-20190328211700-ab21143f2384/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=
golang.org/x/tools v0.0.0-20190425150028-36563e24a262/go.mod h1:RgjU9mgBXZiqYHBnxXauZ1Gv1EHHAz9KjViQ78xBX0Q=
golang.org/x/tools v0.0.0-20190506145303-2d16b83fe98c/go.mod h1:RgjU9mgBXZiqYHBnxXauZ1Gv1EHHAz9KjViQ78xBX0Q=
golang.org/x/tools v0.0.0-20190524140312-2c0ae7006135/go.mod h1:RgjU9mgBXZiqYHBnxXauZ1Gv1EHHAz9KjViQ78xBX0Q=
golang.org/x/tools v0.0.0-20190606124116-d0a3d012864b/go.mod h1:/rFqwRUd4F7ZHNgwSSTFct+R/Kf4OFW1sUzUTQQTgfc=
golang.org/x/tools v0.0.0-20190614205625-5aca471b1d59/go.mod h1:/rFqwRUd4F7ZHNgwSSTFct+R/Kf4OFW1sUzUTQQTgfc=
golang.org/x/tools v0.0.0-20190621195816-6e04913cbbac/go.mod h1:/rFqwRUd4F7ZHNgwSSTFct+R/Kf4OFW1sUzUTQQTgfc=
golang.org/x/tools v0.0.0-20190628153133-6cdbf07be9d0/go.mod h1:/rFqwRUd4F7ZHNgwSSTFct+R/Kf4OFW1sUzUTQQTgfc=
golang.org/x/tools v0.0.0-20190816200558-6889da9d5479/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20190911174233-4f2ddba30aff/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191012152004-8de300cfc20a/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191029041327-9cc4af7d6b2c/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191029190741-b9c20aec41a5/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191112195655-aa38f8e97acc/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191113191852-77e3bb0ad9e7/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191115202509-3a792d9c32b2/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191119224855-298f0cb1881e/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191125144606-a911d9008d1f/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191130070609-6e064ea0cf2d/go.mod h1:b+2E5dAYhXwXZwtnZ6UAqBI28+e2cm9otk0dWdXHAEo=
golang.org/x/tools v0.0.0-20191216173652-a0e659d51361/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20191227053925-7b8e75db28f4/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200117161641-43d50277825c/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200122220014-bf1340f18c4a/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200130002326-2f3ba24bd6e7/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200204074204-1cc6d1ef6c74/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200207183749-b753a1ba74fa/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200212150539-ea181f53ac56/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200224181240-023911ca70b2/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200227222343-706bc42d1f0d/go.mod h1:TB2adYChydJhpapKDTa4BR/hXlZSLoq2Wpct/0txZ28=
golang.org/x/tools v0.0.0-20200304193943-95d2e580d8eb/go.mod h1:o4KQGtdN14AW+yjsvvwRTJJuXz8XRtIHtEnmAXLyFUw=
golang.org/x/tools v0.0.0-20200312045724-11d5b4c81c7d/go.mod h1:o4KQGtdN14AW+yjsvvwRTJJuXz8XRtIHtEnmAXLyFUw=
golang.org/x/tools v0.0.0-20200331025713-a30bf2db82d4/go.mod h1:Sl4aGygMT6LrqrWclx+PTx3U+LnKx/seiNR+3G19Ar8=
golang.org/x/tools v0.0.0-20200501065659-ab2804fb9c9d/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
golang.org/x/tools v0.0.0-20200512131952-2bc93b1c0c88/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
golang.org/x/tools v0.0.0-20200515010526-7d3b6ebf133d/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
golang.org/x/tools v0.0.0-20200618134242-20370b0cb4b2/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
golang.org/x/tools v0.0.0-20200619180055-7c47624df98f/go.mod h1:EkVYQZoAsY45+roYkvgYkIh4xh/qjgUK9TdY2XT94GE=
golang.org/x/tools v0.0.0-20200729194436-6467de6f59a7/go.mod h1:njjCfa9FT2d7l9Bc6FUM5FLjQPp3cFF28FI3qnDFljA=
golang.org/x/tools v0.0.0-20200804011535-6c149bb5ef0d/go.mod h1:njjCfa9FT2d7l9Bc6FUM5FLjQPp3cFF28FI3qnDFljA=
golang.org/x/tools v0.0.0-20200825202427-b303f430e36d/go.mod h1:njjCfa9FT2d7l9Bc6FUM5FLjQPp3cFF28FI3qnDFljA=
golang.org/x/tools v0.0.0-20200904185747-39188db58858/go.mod h1:Cj7w3i3Rnn0Xh82ur9kSqwfTHTeVxaDqrfMjpcNT6bE=
golang.org/x/tools v0.0.0-20201110124207-079ba7bd75cd/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
golang.org/x/tools v0.0.0-20201201161351-ac6f37ff4c2a/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
golang.org/x/tools v0.0.0-20201208233053-a543418bbed2/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
golang.org/x/tools v0.0.0-20210105154028-b0ab187a4818/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
golang.org/x/tools v0.0.0-20210106214847-113979e3529a/go.mod h1:emZCQorbCU4vsT4fOWvOPXz4eW1wZW4PmDk9uLelYpA=
golang.org/x/tools v0.1.0/go.mod h1:xkSsbof2nBLbhDlRMhhhyNLN/zl3eTqcnHD5viDpcZ0=
golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
golang.org/x/xerrors v0.0.0-20191204190536-9bdfabe68543/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 h1:go1bK/D/BFZV2I8cIQd1NKEZ+0owSTG1fDTci4IqFcE=
golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1/go.mod h1:I/5z698sn9Ka8TeJc9MKroUUfqBBauWjQqLJ2OPfmY0=
google.golang.org/api v0.4.0/go.mod h1:8k5glujaEP+g9n7WNsDg8QP6cUVNI86fCNMcbazEtwE=
google.golang.org/api v0.7.0/go.mod h1:WtwebWUNSVBH/HAw79HIFXZNqEvBhG+Ra+ax0hx3E3M=
google.golang.org/api v0.8.0/go.mod h1:o4eAsZoiT+ibD93RtjEohWalFOjRDx6CVaqeizhEnKg=
google.golang.org/api v0.9.0/go.mod h1:o4eAsZoiT+ibD93RtjEohWalFOjRDx6CVaqeizhEnKg=
google.golang.org/api v0.13.0/go.mod h1:iLdEw5Ide6rF15KTC1Kkl0iskquN2gFfn9o9XIsbkAI=
google.golang.org/api v0.14.0/go.mod h1:iLdEw5Ide6rF15KTC1Kkl0iskquN2gFfn9o9XIsbkAI=
google.golang.org/api v0.15.0/go.mod h1:iLdEw5Ide6rF15KTC1Kkl0iskquN2gFfn9o9XIsbkAI=
google.golang.org/api v0.17.0/go.mod h1:BwFmGc8tA3vsd7r/7kR8DY7iEEGSU04BFxCo5jP/sfE=
google.golang.org/api v0.18.0/go.mod h1:BwFmGc8tA3vsd7r/7kR8DY7iEEGSU04BFxCo5jP/sfE=
google.golang.org/api v0.19.0/go.mod h1:BwFmGc8tA3vsd7r/7kR8DY7iEEGSU04BFxCo5jP/sfE=
google.golang.org/api v0.20.0/go.mod h1:BwFmGc8tA3vsd7r/7kR8DY7iEEGSU04BFxCo5jP/sfE=
google.golang.org/api v0.22.0/go.mod h1:BwFmGc8tA3vsd7r/7kR8DY7iEEGSU04BFxCo5jP/sfE=
google.golang.org/api v0.24.0/go.mod h1:lIXQywCXRcnZPGlsd8NbLnOjtAoL6em04bJ9+z0MncE=
google.golang.org/api v0.28.0/go.mod h1:lIXQywCXRcnZPGlsd8NbLnOjtAoL6em04bJ9+z0MncE=
google.golang.org/api v0.29.0/go.mod h1:Lcubydp8VUV7KeIHD9z2Bys/sm/vGKnG1UHuDBSrHWM=
google.golang.org/api v0.30.0/go.mod h1:QGmEvQ87FHZNiUVJkT14jQNYJ4ZJjdRF23ZXz5138Fc=
google.golang.org/api v0.35.0/go.mod h1:/XrVsuzM0rZmrsbjJutiuftIzeuTQcEeaYcSk/mQ1dg=
google.golang.org/api v0.36.0/go.mod h1:+z5ficQTmoYpPn8LCUNVpK5I7hwkpjbcgqA7I34qYtE=
google.golang.org/api v0.40.0/go.mod h1:fYKFpnQN0DsDSKRVRcQSDQNtqWPfM9i+zNPxepjRCQ8=
google.golang.org/api v0.41.0/go.mod h1:RkxM5lITDfTzmyKFPt+wGrCJbVfniCr2ool8kTBzRTU=
google.golang.org/api v0.43.0/go.mod h1:nQsDGjRXMo4lvh5hP0TKqF244gqhGcr/YSIykhUk/94=
google.golang.org/api v0.44.0/go.mod h1:EBOGZqzyhtvMDoxwS97ctnh0zUmYY6CxqXsc1AvkYD8=
google.golang.org/appengine v1.4.0/go.mod h1:xpcJRLb0r/rnEns0DIKYYv+WjYCduHsrkT7/EB5XEv4=
google.golang.org/appengine v1.5.0/go.mod h1:xpcJRLb0r/rnEns0DIKYYv+WjYCduHsrkT7/EB5XEv4=
google.golang.org/appengine v1.6.1/go.mod h1:i06prIuMbXzDqacNJfV5OdTW448YApPu5ww/cMBSeb0=
google.golang.org/appengine v1.6.5/go.mod h1:8WjMMxjGQR8xUklV/ARdw2HLXBOI7O7uCIDZVag1xfc=
google.golang.org/appengine v1.6.6/go.mod h1:8WjMMxjGQR8xUklV/ARdw2HLXBOI7O7uCIDZVag1xfc=
google.golang.org/appengine v1.6.7 h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=
google.golang.org/appengine v1.6.7/go.mod h1:8WjMMxjGQR8xUklV/ARdw2HLXBOI7O7uCIDZVag1xfc=
google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884 h1:fiNLklpBwWK1mth30Hlwk+fcdBmIALlgF5iy77O37Ig=
google.golang.org/genproto v0.0.0-20200513103714-09dca8ec2884/go.mod h1:55QSHmfGQM9UVYDPBsyGGes0y52j32PQ3BqQfXhyH3c=
google.golang.org/grpc v1.26.0 h1:2dTRdpdFEEhJYQD8EMLB61nnrzSCTbG38PhqdhvOltg=
google.golang.org/grpc v1.26.0/go.mod h1:qbnxyOmOxrQa7FizSgH+ReBfzJrCY1pSN7KXBS8abTk=
google.golang.org/protobuf v1.25.0/go.mod h1:9JNX74DMeImyA3h4bdi1ymwjUzf21/xIlbajtzgsN7c=
gopkg.in/alecthomas/kingpin.v2 v2.2.6/go.mod h1:FMv+mEhP44yOT+4EoQTLFTRgOQ1FBLkstjWtayDeSgw=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f h1:BLraFXnmrev5lT+xlilqcH8XK9/i0At2xKjWk4p6zsU=
gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/errgo.v2 v2.1.0/go.mod h1:hNsd1EY+bozCKY1Ytp96fpM3vjJbqLJn88ws8XvfDNI=
gopkg.in/fsnotify.v1 v1.4.7/go.mod h1:Tz8NjZHkW78fSQdbUxIjBTcgA1z1m8ZHf0WmKUhAMys=
gopkg.in/inf.v0 v0.9.1/go.mod h1:cWUDdTG/fYaXco+Dcufb5Vnc6Gp2YChqWtbxRZE0mXw=
gopkg.in/ini.v1 v1.51.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
gopkg.in/ini.v1 v1.62.0 h1:duBzk771uxoUuOlyRLkHsygud9+5lrlGjdFBb4mSKDU=
gopkg.in/ini.v1 v1.62.0/go.mod h1:pNLf8WUiyNEtQjuu5G5vTm06TEv9tsIgeAvK8hOrP4k=
gopkg.in/resty.v1 v1.12.0/go.mod h1:mDo4pnntr5jdWRML875a/NmxYqAlA73dVijT2AXvQQo=
gopkg.in/square/go-jose.v2 v2.5.1 h1:7odma5RETjNHWJnR32wx8t+Io4djHE1PqxCFx3iiZ2w=
gopkg.in/square/go-jose.v2 v2.5.1/go.mod h1:M9dMgbHiYLoDGQrXy7OpJDJWiKiU//h+vD76mk0e1AI=
gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7/go.mod h1:dt/ZhP58zS4L8KSrWDmTeBkI65Dw0HsyUHuEVlX15mw=
gopkg.in/yaml.v2 v2.0.0-20170812160011-eb3733d160e7/go.mod h1:JAlM8MvJe8wmxCU4Bli9HhUf9+ttbYbLASfIpnQbh74=
gopkg.in/yaml.v2 v2.2.1/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
gopkg.in/yaml.v2 v2.2.2/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
gopkg.in/yaml.v2 v2.2.3/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
gopkg.in/yaml.v2 v2.2.4/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
gopkg.in/yaml.v2 v2.2.5/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
gopkg.in/yaml.v2 v2.2.8/go.mod h1:hI93XBmqTisBFMUTm0b8Fm+jr3Dg1NNxqwp+5A1VGuI=
gopkg.in/yaml.v2 v2.4.0 h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=
gopkg.in/yaml.v2 v2.4.0/go.mod h1:RDklbk79AGWmwhnvt/jBztapEOGDOx6ZbXqjP6csGnQ=
gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b h1:h8qDotaEPuJATrMmW04NCwg7v22aHH28wwpauUhK9Oo=
gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
honnef.co/go/tools v0.0.0-20190102054323-c2f93a96b099/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
honnef.co/go/tools v0.0.0-20190106161140-3f1c8253044a/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
honnef.co/go/tools v0.0.0-20190418001031-e561f6794a2a/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
honnef.co/go/tools v0.0.0-20190523083050-ea95bdfd59fc/go.mod h1:rf3lG4BRIbNafJWhAfAdb/ePZxsR/4RtNHQocxwk9r4=
honnef.co/go/tools v0.0.1-2019.2.3/go.mod h1:a3bituU0lyd329TUQxRnasdCoJDkEUEAqEt0JzvZhAg=
honnef.co/go/tools v0.0.1-2020.1.3/go.mod h1:X/FiERA/W4tHapMX5mGpAtMSVEeEUOyHaw9vFzvIQ3k=
honnef.co/go/tools v0.0.1-2020.1.4/go.mod h1:X/FiERA/W4tHapMX5mGpAtMSVEeEUOyHaw9vFzvIQ3k=
k8s.io/apimachinery v0.21.1 h1:Q6XuHGlj2xc+hlMCvqyYfbv3H7SRGn2c8NycxJquDVs=
k8s.io/apimachinery v0.21.1/go.mod h1:jbreFvJo3ov9rj7eWT7+sYiRx+qZuCYXwWT1bcDswPY=
k8s.io/gengo v0.0.0-20200413195148-3a45101e95ac/go.mod h1:ezvh/TsK7cY6rbqRK0oQQ8IAqLxYwwyPxAX1Pzy0ii0=
k8s.io/klog/v2 v2.0.0/go.mod h1:PBfzABfn139FHAV07az/IF9Wp1bkk3vpT2XSJ76fSDE=
k8s.io/klog/v2 v2.8.0 h1:Q3gmuM9hKEjefWFFYF0Mat+YyFJvsUyYuwyNNJ5C9Ts=
k8s.io/klog/v2 v2.8.0/go.mod h1:hy9LJ/NvuK+iVyP4Ehqva4HxZG/oXyIS3n3Jmire4Ec=
k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7/go.mod h1:wXW5VT87nVfh/iLV8FpR2uDvrFyomxbtb1KivDbvPTE=
rsc.io/binaryregexp v0.2.0/go.mod h1:qTv7/COck+e2FymRvadv62gMdZztPaShugOCi3I+8D8=
rsc.io/quote/v3 v3.1.0/go.mod h1:yEA65RcK8LyAZtP9Kv3t0HmxON59tX3rD+tICJqUlj0=
rsc.io/sampler v1.3.0/go.mod h1:T1hPZKmBbMNahiBKFy5HrXp6adAjACjK9JXDnKaTXpA=
sigs.k8s.io/structured-merge-diff/v4 v4.0.2/go.mod h1:bJZC9H9iH24zzfZ/41RGcq60oK1F7G282QMXDPYydCw=
sigs.k8s.io/structured-merge-diff/v4 v4.1.0/go.mod h1:bJZC9H9iH24zzfZ/41RGcq60oK1F7G282QMXDPYydCw=
sigs.k8s.io/yaml v1.2.0 h1:kr/MCeFWJWTwyaHoR9c8EjH9OumOmoF9YGiZd7lFm/Q=
sigs.k8s.io/yaml v1.2.0/go.mod h1:yfXDCHCao9+ENCvLSE62v9VSji2MKu5jeNfTrofGhJc=
`,
	})

	t.files = append(t.files, &templateFile{
		name: "CHANGELOG.md",
		body: `
# Changelog

名称 | 说明
------------|----------
Added       | 添加新功能
Changed     | 功能的变更
Deprecated  | 未来会删除
Removed     | 之前为Deprecated状态，此版本被移除
Fixed       | 功能的修复
Security    | 有关安全问题的修复

## [Unreleased]

### Added

- 已完成的功能，未正式发布

## [0.1.0] - 2020-03-28

### Added

- 首次发布
`,
	})

	t.files = append(t.files, &templateFile{
		name: "Makefile",
		body: `
.PHONY: precheck clean version proto

include scripts/env

# 全局通用变量
GO              := go
GORUN           := ${GO} run
GOPATH          := $(shell ${GO} env GOPATH)
GOARCH          ?= $(shell ${GO} env GOARCH)
GOBUILD         := ${GO} build
WORKDIR         := $(shell pwd)
GOHOSTOS        := $(shell ${GO} env GOHOSTOS)
SHORTNAME       := $(shell basename $(shell pwd))
NAMESPACE       ?= $(shell basename $(shell cd ../;pwd))

# 自动化版本号
GIT_COMMIT      := $(shell git rev-parse HEAD 2>/dev/null)
GIT_BRANCH      := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null)
BUILD_DATE      := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_DATE     := $(shell git --no-pager log -1 --format='%ct' 2>/dev/null)
PREFIX_VERSION  := $(shell ./scripts/version.sh prefix)
RELEASE_VERSION ?= $(shell ./scripts/version.sh release)
BUILD_LD_FLAGS  := "-X 'github.com/grpc-kit/pkg/version.CliVersion=${CLI_VERSION}' \
                -X 'github.com/grpc-kit/pkg/version.GitCommit=${GIT_COMMIT}' \
                -X 'github.com/grpc-kit/pkg/version.GitBranch=${GIT_BRANCH}' \
                -X 'github.com/grpc-kit/pkg/version.BuildDate=${BUILD_DATE}' \
                -X 'github.com/grpc-kit/pkg/version.CommitUnixTime=${COMMIT_DATE}' \
                -X 'github.com/grpc-kit/pkg/version.ReleaseVersion=${RELEASE_VERSION}'"

# 构建Docker容器变量
IMAGE_FROM      ?= scratch
IMAGE_HOST      ?= hub.docker.com
BUILD_GOOS      ?= $(shell ${GO} env GOOS)

precheck:
	@echo ">> precheck environment"
	@./scripts/precheck.sh

run: generate
	@${GORUN} -ldflags ${BUILD_LD_FLAGS} cmd/server/main.go -c config/app-dev-local.yaml

build: clean generate
	@mkdir build
	@mkdir build/deploy
	@GOOS=${BUILD_GOOS} ${GOBUILD} -ldflags ${BUILD_LD_FLAGS} -o build/service cmd/server/main.go

clean:
	@echo ">> clean build"
	@rm -rf build/

generate: precheck
	@echo ">> generation release version"
	@./scripts/version.sh update

	@echo ">> generation assets to static code"
	@${GO} generate api/generate.go >> /dev/null

	@echo ">> generation code from proto files"
	@./scripts/genproto.sh

build-docker: build
	@echo ">> build docker"

	@IMAGE_FROM=${IMAGE_FROM} \
	IMAGE_HOST=${IMAGE_HOST} \
	NAMESPACE=${NAMESPACE} \
	SHORTNAME=${SHORTNAME} \
	RELEASE_VERSION=${RELEASE_VERSION} \
	./scripts/docker.sh
`,
	})
}
