"use strict";(self.webpackChunkmeteor=self.webpackChunkmeteor||[]).push([[934],{3905:(e,t,r)=>{r.d(t,{Zo:()=>c,kt:()=>f});var n=r(7294);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function o(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function l(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?o(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):o(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function i(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var s=n.createContext({}),p=function(e){var t=n.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):l(l({},t),e)),r},c=function(e){var t=p(e.components);return n.createElement(s.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},u=n.forwardRef((function(e,t){var r=e.components,a=e.mdxType,o=e.originalType,s=e.parentName,c=i(e,["components","mdxType","originalType","parentName"]),u=p(r),f=a,d=u["".concat(s,".").concat(f)]||u[f]||m[f]||o;return r?n.createElement(d,l(l({ref:t},c),{},{components:r})):n.createElement(d,l({ref:t},c))}));function f(e,t){var r=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var o=r.length,l=new Array(o);l[0]=u;var i={};for(var s in t)hasOwnProperty.call(t,s)&&(i[s]=t[s]);i.originalType=e,i.mdxType="string"==typeof e?e:a,l[1]=i;for(var p=2;p<o;p++)l[p]=r[p];return n.createElement.apply(null,l)}return n.createElement.apply(null,r)}u.displayName="MDXCreateElement"},9728:(e,t,r)=>{r.r(t),r.d(t,{assets:()=>s,contentTitle:()=>l,default:()=>m,frontMatter:()=>o,metadata:()=>i,toc:()=>p});var n=r(7462),a=(r(7294),r(3905));const o={},l="Processors",i={unversionedId:"reference/processors",id:"reference/processors",title:"Processors",description:"Enrich",source:"@site/docs/reference/processors.md",sourceDirName:"reference",slug:"/reference/processors",permalink:"/meteor/docs/reference/processors",draft:!1,editUrl:"https://github.com/odpf/meteor/edit/master/docs/docs/reference/processors.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Extractors",permalink:"/meteor/docs/reference/extractors"},next:{title:"Sinks",permalink:"/meteor/docs/reference/sinks"}},s={},p=[{value:"Enrich",id:"enrich",level:2},{value:"Configs",id:"configs",level:3},{value:"Sample usage",id:"sample-usage",level:3},{value:"Labels",id:"labels",level:2},{value:"Script",id:"script",level:2}],c={toc:p};function m(e){let{components:t,...r}=e;return(0,a.kt)("wrapper",(0,n.Z)({},c,r,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"processors"},"Processors"),(0,a.kt)("h2",{id:"enrich"},"Enrich"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"enrich")),(0,a.kt)("p",null,"Enrich extra fields to metadata."),(0,a.kt)("h3",{id:"configs"},"Configs"),(0,a.kt)("table",null,(0,a.kt)("thead",{parentName:"table"},(0,a.kt)("tr",{parentName:"thead"},(0,a.kt)("th",{parentName:"tr",align:"left"},"Key"),(0,a.kt)("th",{parentName:"tr",align:"left"},"Value"),(0,a.kt)("th",{parentName:"tr",align:"left"},"Example"),(0,a.kt)("th",{parentName:"tr",align:"left"},"Description"),(0,a.kt)("th",{parentName:"tr",align:"left"}),(0,a.kt)("th",{parentName:"tr",align:"left"}))),(0,a.kt)("tbody",{parentName:"table"},(0,a.kt)("tr",{parentName:"tbody"},(0,a.kt)("td",{parentName:"tr",align:"left"},(0,a.kt)("inlineCode",{parentName:"td"},"{field_name}")),(0,a.kt)("td",{parentName:"tr",align:"left"},"`","string"),(0,a.kt)("td",{parentName:"tr",align:"left"},"number","`"),(0,a.kt)("td",{parentName:"tr",align:"left"},(0,a.kt)("inlineCode",{parentName:"td"},"{field_value}")),(0,a.kt)("td",{parentName:"tr",align:"left"},"Dynamic field and value"),(0,a.kt)("td",{parentName:"tr",align:"left"},(0,a.kt)("em",{parentName:"td"},"required"))))),(0,a.kt)("h3",{id:"sample-usage"},"Sample usage"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre",className:"language-yaml"},"processors:\n  - name: enrich\n    config:\n      fieldA: valueA\n      fieldB: valueB\n")),(0,a.kt)("h2",{id:"labels"},"Labels"),(0,a.kt)("p",null,(0,a.kt)("inlineCode",{parentName:"p"},"labels")),(0,a.kt)("p",null,"This processor will append Asset's Labels with value from given config."),(0,a.kt)("p",null,(0,a.kt)("a",{parentName:"p",href:"https://github.com/odpf/meteor/blob/main/plugins/processors/labels/README.md"},"More details")),(0,a.kt)("h2",{id:"script"},"Script"),(0,a.kt)("p",null,"Script processor uses the user specified script to transform each asset emitted\nfrom the extractor. Currently, ",(0,a.kt)("a",{parentName:"p",href:"https://github.com/d5/tengo"},"Tengo")," is the only supported script\nengine."),(0,a.kt)("p",null,(0,a.kt)("a",{parentName:"p",href:"https://github.com/odpf/meteor/blob/main/plugins/processors/script/README.md"},"More details")))}m.isMDXComponent=!0}}]);