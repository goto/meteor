"use strict";(self.webpackChunkmeteor=self.webpackChunkmeteor||[]).push([[482],{3905:function(e,t,n){n.d(t,{Zo:function(){return p},kt:function(){return m}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function a(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function c(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},i=Object.keys(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var i=Object.getOwnPropertySymbols(e);for(r=0;r<i.length;r++)n=i[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var u=r.createContext({}),s=function(e){var t=r.useContext(u),n=t;return e&&(n="function"==typeof e?e(t):a(a({},t),e)),n},p=function(e){var t=s(e.components);return r.createElement(u.Provider,{value:t},e.children)},l={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,i=e.originalType,u=e.parentName,p=c(e,["components","mdxType","originalType","parentName"]),d=s(n),m=o,f=d["".concat(u,".").concat(m)]||d[m]||l[m]||i;return n?r.createElement(f,a(a({ref:t},p),{},{components:n})):r.createElement(f,a({ref:t},p))}));function m(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var i=n.length,a=new Array(i);a[0]=d;var c={};for(var u in t)hasOwnProperty.call(t,u)&&(c[u]=t[u]);c.originalType=e,c.mdxType="string"==typeof e?e:o,a[1]=c;for(var s=2;s<i;s++)a[s]=n[s];return r.createElement.apply(null,a)}return r.createElement.apply(null,n)}d.displayName="MDXCreateElement"},6096:function(e,t,n){n.r(t),n.d(t,{frontMatter:function(){return c},contentTitle:function(){return u},metadata:function(){return s},toc:function(){return p},default:function(){return d}});var r=n(7462),o=n(3366),i=(n(7294),n(3905)),a=["components"],c={},u="Deployment",s={unversionedId:"guides/deployment",id:"guides/deployment",isDocsHomePage:!1,title:"Deployment",description:"After we are done with running and veryfing that the recipes works with the data-source and sink you have mentioned.",source:"@site/docs/guides/4_deployment.md",sourceDirName:"guides",slug:"/guides/deployment",permalink:"/meteor/docs/guides/deployment",editUrl:"https://github.com/odpf/meteor/edit/master/docs/docs/guides/4_deployment.md",tags:[],version:"current",sidebarPosition:4,frontMatter:{},sidebar:"docsSidebar",previous:{title:"Running Meteor",permalink:"/meteor/docs/guides/run_recipes"},next:{title:"Concepts",permalink:"/meteor/docs/concepts/overview"}},p=[],l={toc:p};function d(e){var t=e.components,n=(0,o.Z)(e,a);return(0,i.kt)("wrapper",(0,r.Z)({},l,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"deployment"},"Deployment"),(0,i.kt)("p",null,"After we are done with running and veryfing that the recipes works with the data-source and sink you have mentioned.\nYou may want to automate the process of metadata collection on some regular basis as a cron job.\nOne can setup user for the same."),(0,i.kt)("p",null,"In ODPF we use helm chart to set it up, and you can refer the same ",(0,i.kt)("a",{parentName:"p",href:"https://github.com/odpf/charts"},"here"),"."))}d.isMDXComponent=!0}}]);