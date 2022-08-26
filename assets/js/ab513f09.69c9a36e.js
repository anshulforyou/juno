"use strict";(self.webpackChunkjuno_docs=self.webpackChunkjuno_docs||[]).push([[1538],{6074:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>a,default:()=>s,frontMatter:()=>i,metadata:()=>l,toc:()=>u});var r=t(7462),o=(t(7294),t(3905));const i={title:"Configuration File",position:2},a=void 0,l={unversionedId:"running/config",id:"running/config",title:"Configuration File",description:"Juno uses yaml format to for its configuration file. It can be provided via the following flag:",source:"@site/docs/running/config.mdx",sourceDirName:"running",slug:"/running/config",permalink:"/docs/running/config",draft:!1,editUrl:"https://github.com/NethermindEth/juno/tree/main/docs/docs/running/config.mdx",tags:[],version:"current",frontMatter:{title:"Configuration File",position:2},sidebar:"tutorialSidebar",previous:{title:"Running",permalink:"/docs/category/running"},next:{title:"Docker Execution",permalink:"/docs/running/docker"}},c={},u=[],p={toc:u};function s(e){let{components:n,...t}=e;return(0,o.kt)("wrapper",(0,r.Z)({},p,t,{components:n,mdxType:"MDXLayout"}),(0,o.kt)("p",null,"Juno uses ",(0,o.kt)("inlineCode",{parentName:"p"},"yaml")," format to for its configuration file. It can be provided via the following flag:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},'$ ./build/juno --config "path/to/config/file"\n')),(0,o.kt)("h1",{id:"junoyaml"},"juno.yaml"),(0,o.kt)("p",null,"Each of the configuration options has a corresponding CLI flag. For example, the CLI flag for\nthe database path is ",(0,o.kt)("inlineCode",{parentName:"p"},"--db-path"),"."),(0,o.kt)("p",null,"To see a description of the configuration options, run ",(0,o.kt)("inlineCode",{parentName:"p"},"./build/juno --help"),"."),(0,o.kt)("p",null,"Example ",(0,o.kt)("inlineCode",{parentName:"p"},"juno.yaml"),":"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-yaml"},'verbosity: "debug"\nrpc-port: 4576\nmetrics: true\ndb-path: "/home/.juno"\nnetwork: 1\neth-node: "https://some-ethnode:5673"\n')))}s.isMDXComponent=!0},3905:(e,n,t)=>{t.d(n,{Zo:()=>p,kt:()=>d});var r=t(7294);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function i(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function a(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?i(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):i(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function l(e,n){if(null==e)return{};var t,r,o=function(e,n){if(null==e)return{};var t,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)t=i[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var c=r.createContext({}),u=function(e){var n=r.useContext(c),t=n;return e&&(t="function"==typeof e?e(n):a(a({},n),e)),t},p=function(e){var n=u(e.components);return r.createElement(c.Provider,{value:n},e.children)},s={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},f=r.forwardRef((function(e,n){var t=e.components,o=e.mdxType,i=e.originalType,c=e.parentName,p=l(e,["components","mdxType","originalType","parentName"]),f=u(t),d=o,m=f["".concat(c,".").concat(d)]||f[d]||s[d]||i;return t?r.createElement(m,a(a({ref:n},p),{},{components:t})):r.createElement(m,a({ref:n},p))}));function d(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var i=t.length,a=new Array(i);a[0]=f;var l={};for(var c in n)hasOwnProperty.call(n,c)&&(l[c]=n[c]);l.originalType=e,l.mdxType="string"==typeof e?e:o,a[1]=l;for(var u=2;u<i;u++)a[u]=t[u];return r.createElement.apply(null,a)}return r.createElement.apply(null,t)}f.displayName="MDXCreateElement"}}]);