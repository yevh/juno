"use strict";(self.webpackChunkjuno_docs=self.webpackChunkjuno_docs||[]).push([[6915],{549:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>u,contentTitle:()=>i,default:()=>p,frontMatter:()=>o,metadata:()=>c,toc:()=>d});var r=t(4848),a=t(8453),s=t(3859),l=t(9365);const o={title:"JSON-RPC Interface"},i="JSON-RPC Interface :globe_with_meridians:",c={id:"json-rpc",title:"JSON-RPC Interface",description:"Interacting with Juno requires sending requests to specific JSON-RPC API methods. Juno supports all of Starknet's Node API Endpoints over HTTP and WebSocket.",source:"@site/versioned_docs/version-0.12.0/json-rpc.md",sourceDirName:".",slug:"/json-rpc",permalink:"/json-rpc",draft:!1,unlisted:!1,tags:[],version:"0.12.0",frontMatter:{title:"JSON-RPC Interface"},sidebar:"main",previous:{title:"Updating Juno",permalink:"/updating"},next:{title:"WebSocket Interface",permalink:"/websocket"}},u={},d=[{value:"Enable the JSON-RPC server",id:"enable-the-json-rpc-server",level:2},{value:"Making JSON-RPC requests",id:"making-json-rpc-requests",level:2},{value:"Supported Starknet API versions",id:"supported-starknet-api-versions",level:2}];function h(e){const n={a:"a",code:"code",h1:"h1",h2:"h2",li:"li",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,a.R)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsxs)(n.h1,{id:"json-rpc-interface-globe_with_meridians",children:["JSON-RPC Interface ","\ud83c\udf10"]}),"\n",(0,r.jsxs)(n.p,{children:["Interacting with Juno requires sending requests to specific JSON-RPC API methods. Juno supports all of ",(0,r.jsx)(n.a,{href:"https://playground.open-rpc.org/?uiSchema%5BappBar%5D%5Bui:splitView%5D=false&schemaUrl=https://raw.githubusercontent.com/starkware-libs/starknet-specs/v0.7.0/api/starknet_api_openrpc.json&uiSchema%5BappBar%5D%5Bui:input%5D=false&uiSchema%5BappBar%5D%5Bui:darkMode%5D=true&uiSchema%5BappBar%5D%5Bui:examplesDropdown%5D=false",children:"Starknet's Node API Endpoints"})," over HTTP and ",(0,r.jsx)(n.a,{href:"websocket",children:"WebSocket"}),"."]}),"\n",(0,r.jsx)(n.h2,{id:"enable-the-json-rpc-server",children:"Enable the JSON-RPC server"}),"\n",(0,r.jsx)(n.p,{children:"To enable the JSON-RPC interface, use the following configuration options:"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"http"}),": Enables the HTTP RPC server on the default port and interface (disabled by default)."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"http-host"}),": The interface on which the HTTP RPC server will listen for requests. If skipped, it defaults to ",(0,r.jsx)(n.code,{children:"localhost"}),"."]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.code,{children:"http-port"}),": The port on which the HTTP server will listen for requests. If skipped, it defaults to ",(0,r.jsx)(n.code,{children:"6060"}),"."]}),"\n"]}),"\n",(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:"# Docker container\ndocker run -d \\\n  --name juno \\\n  -p 6060:6060 \\\n  nethermind/juno \\\n  --http \\\n  --http-port 6060 \\\n  --http-host 0.0.0.0\n\n# Standalone binary\n./build/juno --http --http-port 6060 --http-host 0.0.0.0\n"})}),"\n",(0,r.jsx)(n.h2,{id:"making-json-rpc-requests",children:"Making JSON-RPC requests"}),"\n",(0,r.jsxs)(n.p,{children:["You can use any of ",(0,r.jsx)(n.a,{href:"https://playground.open-rpc.org/?uiSchema%5BappBar%5D%5Bui:splitView%5D=false&schemaUrl=https://raw.githubusercontent.com/starkware-libs/starknet-specs/v0.7.0/api/starknet_api_openrpc.json&uiSchema%5BappBar%5D%5Bui:input%5D=false&uiSchema%5BappBar%5D%5Bui:darkMode%5D=true&uiSchema%5BappBar%5D%5Bui:examplesDropdown%5D=false",children:"Starknet's Node API Endpoints"})," with Juno. Check the availability of Juno with the ",(0,r.jsx)(n.code,{children:"juno_version"})," method:"]}),"\n","\n",(0,r.jsxs)(s.A,{children:[(0,r.jsx)(l.A,{value:"raw",label:"Raw",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-json",children:'{\n  "jsonrpc": "2.0",\n  "method": "juno_version",\n  "params": [],\n  "id": 1\n}\n'})})}),(0,r.jsx)(l.A,{value:"curl",label:"cURL",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:'curl --location \'http://localhost:6060\' \\\n--header \'Content-Type: application/json\' \\\n--data \'{\n    "jsonrpc": "2.0",\n    "method": "juno_version",\n    "params": [],\n    "id": 1\n}\'\n'})})}),(0,r.jsx)(l.A,{value:"response",label:"Response",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-json",children:'{\n  "jsonrpc": "2.0",\n  "result": "v0.11.7",\n  "id": 1\n}\n'})})})]}),"\n",(0,r.jsxs)(n.p,{children:["Get the most recent accepted block hash and number with the ",(0,r.jsx)(n.code,{children:"starknet_blockHashAndNumber"})," method:"]}),"\n",(0,r.jsxs)(s.A,{children:[(0,r.jsx)(l.A,{value:"raw",label:"Raw",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-json",children:'{\n  "jsonrpc": "2.0",\n  "method": "starknet_blockHashAndNumber",\n  "params": [],\n  "id": 1\n}\n'})})}),(0,r.jsx)(l.A,{value:"curl",label:"cURL",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:'curl --location \'http://localhost:6060\' \\\n--header \'Content-Type: application/json\' \\\n--data \'{\n    "jsonrpc": "2.0",\n    "method": "starknet_blockHashAndNumber",\n    "params": [],\n    "id": 1\n}\'\n'})})}),(0,r.jsx)(l.A,{value:"starknetjs",label:"Starknet.js",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-js",children:'const { RpcProvider } = require("starknet");\n\nconst provider = new RpcProvider({\n  nodeUrl: "http://localhost:6060",\n});\n\nprovider.getBlockLatestAccepted().then((blockHashAndNumber) => {\n  console.log(blockHashAndNumber);\n});\n'})})}),(0,r.jsx)(l.A,{value:"starknetgo",label:"Starknet.go",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-go",children:'package main\n\nimport (\n\t"context"\n\t"fmt"\n\t"log"\n\t"github.com/NethermindEth/juno/core/felt"\n\t"github.com/NethermindEth/starknet.go/rpc"\n\t"github.com/NethermindEth/starknet.go/utils"\n)\n\nfunc main() {\n\trpcUrl := "http://localhost:6060"\n\tclient, err := rpc.NewClient(rpcUrl)\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\n\tprovider := rpc.NewProvider(client)\n\tresult, err := provider.BlockHashAndNumber(context.Background())\n\tif err != nil {\n\t\tlog.Fatal(err)\n\t}\n\tfmt.Println("BlockHashAndNumber:", result)\n}\n'})})}),(0,r.jsx)(l.A,{value:"starknetrs",label:"Starknet.rs",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-rust",children:'use starknet::providers::{\n    jsonrpc::{HttpTransport, JsonRpcClient},\n    Provider, Url,\n};\n\n#[tokio::main]\nasync fn main() {\n    let provider = JsonRpcClient::new(HttpTransport::new(\n        Url::parse("http://localhost:6060").unwrap(),\n    ));\n\n    let result = provider.block_hash_and_number().await;\n    match result {\n        Ok(block_hash_and_number) => {\n            println!("{block_hash_and_number:#?}");\n        }\n        Err(err) => {\n            eprintln!("Error: {err}");\n        }\n    }\n}\n'})})}),(0,r.jsx)(l.A,{value:"response",label:"Response",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-json",children:'{\n  "jsonrpc": "2.0",\n  "result": {\n    "block_hash": "0x637ae4d7468bb603c2f16ba7f9118d58c7d7c98a8210260372e83e7c9df443a",\n    "block_number": 640827\n  },\n  "id": 1\n}\n'})})})]}),"\n",(0,r.jsx)(n.h2,{id:"supported-starknet-api-versions",children:"Supported Starknet API versions"}),"\n",(0,r.jsx)(n.p,{children:"Juno supports the following Starknet API versions:"}),"\n",(0,r.jsxs)(n.ul,{children:["\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"v0.7.0"}),": Accessible via endpoints ",(0,r.jsx)(n.code,{children:"/v0_7"}),", ",(0,r.jsx)(n.code,{children:"/rpc/v0_7"}),", or the default ",(0,r.jsx)(n.code,{children:"/"})]}),"\n",(0,r.jsxs)(n.li,{children:[(0,r.jsx)(n.strong,{children:"v0.6.0"}),": Accessible via endpoints ",(0,r.jsx)(n.code,{children:"/v0_6"})," or ",(0,r.jsx)(n.code,{children:"/rpc/v0_6"})]}),"\n"]}),"\n",(0,r.jsx)(n.p,{children:"To use a specific API version, specify the version endpoint in your RPC calls:"}),"\n",(0,r.jsxs)(s.A,{children:[(0,r.jsx)(l.A,{value:"latest",label:"Latest",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:'curl --location \'http://localhost:6060\' \\\n--header \'Content-Type: application/json\' \\\n--data \'{\n    "jsonrpc": "2.0",\n    "method": "starknet_chainId",\n    "params": [],\n    "id": 1\n}\'\n'})})}),(0,r.jsx)(l.A,{value:"v7",label:"v0.7.0",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:'curl --location \'http://localhost:6060/v0_7\' \\\n--header \'Content-Type: application/json\' \\\n--data \'{\n    "jsonrpc": "2.0",\n    "method": "starknet_chainId",\n    "params": [],\n    "id": 1\n}\'\n'})})}),(0,r.jsx)(l.A,{value:"v6",label:"v0.6.0",children:(0,r.jsx)(n.pre,{children:(0,r.jsx)(n.code,{className:"language-bash",children:'curl --location \'http://localhost:6060/v0_6\' \\\n--header \'Content-Type: application/json\' \\\n--data \'{\n    "jsonrpc": "2.0",\n    "method": "starknet_chainId",\n    "params": [],\n    "id": 1\n}\'\n'})})})]})]})}function p(e={}){const{wrapper:n}={...(0,a.R)(),...e.components};return n?(0,r.jsx)(n,{...e,children:(0,r.jsx)(h,{...e})}):h(e)}},9365:(e,n,t)=>{t.d(n,{A:()=>l});t(6540);var r=t(4164);const a={tabItem:"tabItem_Ymn6"};var s=t(4848);function l(e){let{children:n,hidden:t,className:l}=e;return(0,s.jsx)("div",{role:"tabpanel",className:(0,r.A)(a.tabItem,l),hidden:t,children:n})}},3859:(e,n,t)=>{t.d(n,{A:()=>w});var r=t(6540),a=t(4164),s=t(6641),l=t(6347),o=t(205),i=t(8874),c=t(4035),u=t(2993);function d(e){return r.Children.toArray(e).filter((e=>"\n"!==e)).map((e=>{if(!e||(0,r.isValidElement)(e)&&function(e){const{props:n}=e;return!!n&&"object"==typeof n&&"value"in n}(e))return e;throw new Error(`Docusaurus error: Bad <Tabs> child <${"string"==typeof e.type?e.type:e.type.name}>: all children of the <Tabs> component should be <TabItem>, and every <TabItem> should have a unique "value" prop.`)}))?.filter(Boolean)??[]}function h(e){const{values:n,children:t}=e;return(0,r.useMemo)((()=>{const e=n??function(e){return d(e).map((e=>{let{props:{value:n,label:t,attributes:r,default:a}}=e;return{value:n,label:t,attributes:r,default:a}}))}(t);return function(e){const n=(0,c.X)(e,((e,n)=>e.value===n.value));if(n.length>0)throw new Error(`Docusaurus error: Duplicate values "${n.map((e=>e.value)).join(", ")}" found in <Tabs>. Every value needs to be unique.`)}(e),e}),[n,t])}function p(e){let{value:n,tabValues:t}=e;return t.some((e=>e.value===n))}function m(e){let{queryString:n=!1,groupId:t}=e;const a=(0,l.W6)(),s=function(e){let{queryString:n=!1,groupId:t}=e;if("string"==typeof n)return n;if(!1===n)return null;if(!0===n&&!t)throw new Error('Docusaurus error: The <Tabs> component groupId prop is required if queryString=true, because this value is used as the search param name. You can also provide an explicit value such as queryString="my-search-param".');return t??null}({queryString:n,groupId:t});return[(0,i.aZ)(s),(0,r.useCallback)((e=>{if(!s)return;const n=new URLSearchParams(a.location.search);n.set(s,e),a.replace({...a.location,search:n.toString()})}),[s,a])]}function b(e){const{defaultValue:n,queryString:t=!1,groupId:a}=e,s=h(e),[l,i]=(0,r.useState)((()=>function(e){let{defaultValue:n,tabValues:t}=e;if(0===t.length)throw new Error("Docusaurus error: the <Tabs> component requires at least one <TabItem> children component");if(n){if(!p({value:n,tabValues:t}))throw new Error(`Docusaurus error: The <Tabs> has a defaultValue "${n}" but none of its children has the corresponding value. Available values are: ${t.map((e=>e.value)).join(", ")}. If you intend to show no default tab, use defaultValue={null} instead.`);return n}const r=t.find((e=>e.default))??t[0];if(!r)throw new Error("Unexpected error: 0 tabValues");return r.value}({defaultValue:n,tabValues:s}))),[c,d]=m({queryString:t,groupId:a}),[b,j]=function(e){let{groupId:n}=e;const t=function(e){return e?`docusaurus.tab.${e}`:null}(n),[a,s]=(0,u.Dv)(t);return[a,(0,r.useCallback)((e=>{t&&s.set(e)}),[t,s])]}({groupId:a}),f=(()=>{const e=c??b;return p({value:e,tabValues:s})?e:null})();(0,o.A)((()=>{f&&i(f)}),[f]);return{selectedValue:l,selectValue:(0,r.useCallback)((e=>{if(!p({value:e,tabValues:s}))throw new Error(`Can't select invalid tab value=${e}`);i(e),d(e),j(e)}),[d,j,s]),tabValues:s}}var j=t(2303);const f={tabList:"tabList__CuJ",tabItem:"tabItem_LNqP"};var v=t(4848);function x(e){let{className:n,block:t,selectedValue:r,selectValue:l,tabValues:o}=e;const i=[],{blockElementScrollPositionUntilNextRender:c}=(0,s.a_)(),u=e=>{const n=e.currentTarget,t=i.indexOf(n),a=o[t].value;a!==r&&(c(n),l(a))},d=e=>{let n=null;switch(e.key){case"Enter":u(e);break;case"ArrowRight":{const t=i.indexOf(e.currentTarget)+1;n=i[t]??i[0];break}case"ArrowLeft":{const t=i.indexOf(e.currentTarget)-1;n=i[t]??i[i.length-1];break}}n?.focus()};return(0,v.jsx)("ul",{role:"tablist","aria-orientation":"horizontal",className:(0,a.A)("tabs",{"tabs--block":t},n),children:o.map((e=>{let{value:n,label:t,attributes:s}=e;return(0,v.jsx)("li",{role:"tab",tabIndex:r===n?0:-1,"aria-selected":r===n,ref:e=>i.push(e),onKeyDown:d,onClick:u,...s,className:(0,a.A)("tabs__item",f.tabItem,s?.className,{"tabs__item--active":r===n}),children:t??n},n)}))})}function g(e){let{lazy:n,children:t,selectedValue:a}=e;const s=(Array.isArray(t)?t:[t]).filter(Boolean);if(n){const e=s.find((e=>e.props.value===a));return e?(0,r.cloneElement)(e,{className:"margin-top--md"}):null}return(0,v.jsx)("div",{className:"margin-top--md",children:s.map(((e,n)=>(0,r.cloneElement)(e,{key:n,hidden:e.props.value!==a})))})}function k(e){const n=b(e);return(0,v.jsxs)("div",{className:(0,a.A)("tabs-container",f.tabList),children:[(0,v.jsx)(x,{...e,...n}),(0,v.jsx)(g,{...e,...n})]})}function w(e){const n=(0,j.A)();return(0,v.jsx)(k,{...e,children:d(e.children)},String(n))}},8453:(e,n,t)=>{t.d(n,{R:()=>l,x:()=>o});var r=t(6540);const a={},s=r.createContext(a);function l(e){const n=r.useContext(s);return r.useMemo((function(){return"function"==typeof e?e(n):{...n,...e}}),[n,e])}function o(e){let n;return n=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:l(e.components),r.createElement(s.Provider,{value:n},e.children)}}}]);