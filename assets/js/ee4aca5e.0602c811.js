"use strict";(self.webpackChunkjuno_docs=self.webpackChunkjuno_docs||[]).push([[3805],{5132:(e,n,t)=>{t.r(n),t.d(n,{assets:()=>c,contentTitle:()=>l,default:()=>p,frontMatter:()=>a,metadata:()=>u,toc:()=>s});var r=t(7462),o=(t(7294),t(3905));const a={title:"Useful Commands"},l="\ud83d\ude9a Run Locally",u={unversionedId:"running/useful-commands",id:"running/useful-commands",title:"Useful Commands",description:"Install dependencies:",source:"@site/docs/running/useful-commands.mdx",sourceDirName:"running",slug:"/running/useful-commands",permalink:"/docs/running/useful-commands",draft:!1,editUrl:"https://github.com/NethermindEth/juno/tree/main/docs/docs/running/useful-commands.mdx",tags:[],version:"current",frontMatter:{title:"Useful Commands"},sidebar:"tutorialSidebar",previous:{title:"How to Deploy",permalink:"/docs/running/how-to-deploy"},next:{title:"Features",permalink:"/docs/category/features"}},c={},s=[],i={toc:s};function p(e){let{components:n,...t}=e;return(0,o.kt)("wrapper",(0,r.Z)({},i,t,{components:n,mdxType:"MDXLayout"}),(0,o.kt)("h1",{id:"-run-locally"},"\ud83d\ude9a Run Locally"),(0,o.kt)("p",null,"Install dependencies:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},"make install-deps\n")),(0,o.kt)("p",null,"In the ",(0,o.kt)("inlineCode",{parentName:"p"},"juno")," directory, run the following to compile Juno:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},"make juno\n")),(0,o.kt)("p",null,"To run tests:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},"make test\n")),(0,o.kt)("p",null,"For more details of the makefile, just run:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},"make help\n")),(0,o.kt)("p",null,"For a full list of options, use:"),(0,o.kt)("pre",null,(0,o.kt)("code",{parentName:"pre",className:"language-bash"},"./build/juno -h\n")))}p.isMDXComponent=!0},3905:(e,n,t)=>{t.d(n,{Zo:()=>i,kt:()=>d});var r=t(7294);function o(e,n,t){return n in e?Object.defineProperty(e,n,{value:t,enumerable:!0,configurable:!0,writable:!0}):e[n]=t,e}function a(e,n){var t=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);n&&(r=r.filter((function(n){return Object.getOwnPropertyDescriptor(e,n).enumerable}))),t.push.apply(t,r)}return t}function l(e){for(var n=1;n<arguments.length;n++){var t=null!=arguments[n]?arguments[n]:{};n%2?a(Object(t),!0).forEach((function(n){o(e,n,t[n])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(t)):a(Object(t)).forEach((function(n){Object.defineProperty(e,n,Object.getOwnPropertyDescriptor(t,n))}))}return e}function u(e,n){if(null==e)return{};var t,r,o=function(e,n){if(null==e)return{};var t,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)t=a[r],n.indexOf(t)>=0||(o[t]=e[t]);return o}(e,n);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)t=a[r],n.indexOf(t)>=0||Object.prototype.propertyIsEnumerable.call(e,t)&&(o[t]=e[t])}return o}var c=r.createContext({}),s=function(e){var n=r.useContext(c),t=n;return e&&(t="function"==typeof e?e(n):l(l({},n),e)),t},i=function(e){var n=s(e.components);return r.createElement(c.Provider,{value:n},e.children)},p={inlineCode:"code",wrapper:function(e){var n=e.children;return r.createElement(r.Fragment,{},n)}},m=r.forwardRef((function(e,n){var t=e.components,o=e.mdxType,a=e.originalType,c=e.parentName,i=u(e,["components","mdxType","originalType","parentName"]),m=s(t),d=o,f=m["".concat(c,".").concat(d)]||m[d]||p[d]||a;return t?r.createElement(f,l(l({ref:n},i),{},{components:t})):r.createElement(f,l({ref:n},i))}));function d(e,n){var t=arguments,o=n&&n.mdxType;if("string"==typeof e||o){var a=t.length,l=new Array(a);l[0]=m;var u={};for(var c in n)hasOwnProperty.call(n,c)&&(u[c]=n[c]);u.originalType=e,u.mdxType="string"==typeof e?e:o,l[1]=u;for(var s=2;s<a;s++)l[s]=t[s];return r.createElement.apply(null,l)}return r.createElement.apply(null,t)}m.displayName="MDXCreateElement"}}]);