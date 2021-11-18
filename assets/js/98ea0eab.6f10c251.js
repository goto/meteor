"use strict";(self.webpackChunkmeteor=self.webpackChunkmeteor||[]).push([[142],{3905:function(e,t,r){r.d(t,{Zo:function(){return u},kt:function(){return f}});var n=r(7294);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function i(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function a(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?i(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):i(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function c(e,t){if(null==e)return{};var r,n,o=function(e,t){if(null==e)return{};var r,n,o={},i=Object.keys(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||(o[r]=e[r]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(n=0;n<i.length;n++)r=i[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(o[r]=e[r])}return o}var s=n.createContext({}),l=function(e){var t=n.useContext(s),r=t;return e&&(r="function"==typeof e?e(t):a(a({},t),e)),r},u=function(e){var t=l(e.components);return n.createElement(s.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},p=n.forwardRef((function(e,t){var r=e.components,o=e.mdxType,i=e.originalType,s=e.parentName,u=c(e,["components","mdxType","originalType","parentName"]),p=l(r),f=o,m=p["".concat(s,".").concat(f)]||p[f]||d[f]||i;return r?n.createElement(m,a(a({ref:t},u),{},{components:r})):n.createElement(m,a({ref:t},u))}));function f(e,t){var r=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=r.length,a=new Array(i);a[0]=p;var c={};for(var s in t)hasOwnProperty.call(t,s)&&(c[s]=t[s]);c.originalType=e,c.mdxType="string"==typeof e?e:o,a[1]=c;for(var l=2;l<i;l++)a[l]=r[l];return n.createElement.apply(null,a)}return n.createElement.apply(null,r)}p.displayName="MDXCreateElement"},9139:function(e,t,r){r.r(t),r.d(t,{frontMatter:function(){return c},contentTitle:function(){return s},metadata:function(){return l},toc:function(){return u},default:function(){return p}});var n=r(7462),o=r(3366),i=(r(7294),r(3905)),a=["components"],c={},s="Guide",l={unversionedId:"contribute/guide",id:"contribute/guide",isDocsHomePage:!1,title:"Guide",description:"Adding a new Extractor",source:"@site/docs/contribute/guide.md",sourceDirName:"contribute",slug:"/contribute/guide",permalink:"/meteor/docs/contribute/guide",editUrl:"https://github.com/odpf/meteor/edit/master/docs/docs/contribute/guide.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Sinks",permalink:"/meteor/docs/reference/sinks"},next:{title:"Contribution Process",permalink:"/meteor/docs/contribute/contributing"}},u=[{value:"Adding a new Extractor",id:"adding-a-new-extractor",children:[]},{value:"Adding a new Processor",id:"adding-a-new-processor",children:[]},{value:"Adding a new Sink",id:"adding-a-new-sink",children:[]}],d={toc:u};function p(e){var t=e.components,r=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,n.Z)({},d,r,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"guide"},"Guide"),(0,i.kt)("h2",{id:"adding-a-new-extractor"},"Adding a new Extractor"),(0,i.kt)("p",null,"Please follow this list when adding a new Extractor:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Your extractor has to implement one of the ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/core/interface.go"},"defined interfaces"),"."),(0,i.kt)("li",{parentName:"ul"},"Create unit test for the new extractor."),(0,i.kt)("li",{parentName:"ul"},"Add ",(0,i.kt)("a",{parentName:"li",href:"https://pkg.go.dev/go/build#hdr-Build_Constraints"},"build tags")," ",(0,i.kt)("inlineCode",{parentName:"li"},"//+build integration")," on top of your unit test file as shown ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/extractors/mysql/extractor_test.go"},"here"),". This would make sure the test will not be run on when we are testing all unit tests."),(0,i.kt)("li",{parentName:"ul"},"If the source instance is required for testing, Meteor provides a utility to easily create a docker container to help with your test as shown ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/extractors/mysql/extractor_test.go#L35"},"here"),"."),(0,i.kt)("li",{parentName:"ul"},"Register your extractor ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/extractors/populate.go"},"here"),". This is also where you would inject any dependencies needed for your extractor."),(0,i.kt)("li",{parentName:"ul"},"Create a markdown with your extractor details. ","(",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/extractors/mysql/README.md"},"example"),")"),(0,i.kt)("li",{parentName:"ul"},"Add your extractor to one of the extractor list in ",(0,i.kt)("inlineCode",{parentName:"li"},"docs/reference/extractors.md"),"."),(0,i.kt)("li",{parentName:"ul"},"Your extractor should return one of these ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/proto/odpf/meta/data_models.md"},"data models")," as output.")),(0,i.kt)("h2",{id:"adding-a-new-processor"},"Adding a new Processor"),(0,i.kt)("p",null,"Please follow this list when adding a new Processor:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Create unit test for the new processor."),(0,i.kt)("li",{parentName:"ul"},"If the source instance is required for testing, Meteor provides a utility to easily create a docker container to help with your test as shown ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/extractors/mysql/extractor_test.go#L35"},"here"),"."),(0,i.kt)("li",{parentName:"ul"},"Register your processor ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/processors/populate.go"},"here"),". This is also where you would inject any dependencies needed for your processor."),(0,i.kt)("li",{parentName:"ul"},"Update ",(0,i.kt)("inlineCode",{parentName:"li"},"docs/reference/processors.md")," with guide to use the new processor.")),(0,i.kt)("h2",{id:"adding-a-new-sink"},"Adding a new Sink"),(0,i.kt)("p",null,"Please follow this list when adding a new Sink:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},"Create unit test for the new processor."),(0,i.kt)("li",{parentName:"ul"},"If the source instance is required for testing, Meteor provides a utility to easily create a docker container to help with your test as shown ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/extractors/mysql/extractor_test.go#L35"},"here"),"."),(0,i.kt)("li",{parentName:"ul"},"Register your sink ",(0,i.kt)("a",{parentName:"li",href:"https://github.com/odpf/meteor/tree/27f39fe2f83b657d4ecb9eb2c2a8794c6c0671b6/plugins/sinks/populate.go"},"here"),". This is also where you would inject any dependencies needed for your sink."),(0,i.kt)("li",{parentName:"ul"},"Update ",(0,i.kt)("inlineCode",{parentName:"li"},"docs/reference/sinks.md")," with guide to use the new sink.")))}p.isMDXComponent=!0}}]);