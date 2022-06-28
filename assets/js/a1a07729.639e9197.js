"use strict";(self.webpackChunkmeteor=self.webpackChunkmeteor||[]).push([[982],{3905:(e,t,n)=>{n.d(t,{Zo:()=>c,kt:()=>d});var r=n(7294);function i(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function a(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function o(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?a(Object(n),!0).forEach((function(t){i(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):a(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,i=function(e,t){if(null==e)return{};var n,r,i={},a=Object.keys(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||(i[n]=e[n]);return i}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(i[n]=e[n])}return i}var s=r.createContext({}),p=function(e){var t=r.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):o(o({},t),e)),n},c=function(e){var t=p(e.components);return r.createElement(s.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},u=r.forwardRef((function(e,t){var n=e.components,i=e.mdxType,a=e.originalType,s=e.parentName,c=l(e,["components","mdxType","originalType","parentName"]),u=p(n),d=i,g=u["".concat(s,".").concat(d)]||u[d]||m[d]||a;return n?r.createElement(g,o(o({ref:t},c),{},{components:n})):r.createElement(g,o({ref:t},c))}));function d(e,t){var n=arguments,i=t&&t.mdxType;if("string"==typeof e||i){var a=n.length,o=new Array(a);o[0]=u;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:i,o[1]=l;for(var p=2;p<a;p++)o[p]=n[p];return r.createElement.apply(null,o)}return r.createElement.apply(null,n)}u.displayName="MDXCreateElement"},8364:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>s,contentTitle:()=>o,default:()=>m,frontMatter:()=>a,metadata:()=>l,toc:()=>p});var r=n(7462),i=(n(7294),n(3905));const a={},o="Commands",l={unversionedId:"reference/commands",id:"reference/commands",title:"Commands",description:"Meteor currently supports the following commands and these can be utilised after the installation:",source:"@site/docs/reference/commands.md",sourceDirName:"reference",slug:"/reference/commands",permalink:"/meteor/docs/reference/commands",draft:!1,editUrl:"https://github.com/odpf/meteor/edit/master/docs/docs/reference/commands.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Sink",permalink:"/meteor/docs/concepts/sink"},next:{title:"Configuration",permalink:"/meteor/docs/reference/configuration"}},s={},p=[{value:"Listing all the plugins",id:"listing-all-the-plugins",level:2},{value:"Getting Information about plugins",id:"getting-information-about-plugins",level:2},{value:"Generating Sample recipe(s)",id:"generating-sample-recipes",level:2},{value:"Linting recipes",id:"linting-recipes",level:2},{value:"Running recipes",id:"running-recipes",level:2},{value:"get help on commands when stuck",id:"get-help-on-commands-when-stuck",level:2}],c={toc:p};function m(e){let{components:t,...n}=e;return(0,i.kt)("wrapper",(0,r.Z)({},c,n,{components:t,mdxType:"MDXLayout"}),(0,i.kt)("h1",{id:"commands"},"Commands"),(0,i.kt)("p",null,"Meteor currently supports the following commands and these can be utilised after the installation:"),(0,i.kt)("ul",null,(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},"completion: generate the auto completion script for the specified shell")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"#creating-sample-recipes"},"gen"),": The recipe will be printed on standard output.\nSpecify recipe name with the first argument without extension.\nUse comma to separate multiple sinks and processors.")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"#get-help-on-commands-when-stuck"},"help"),": to help the user with meteor.")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"#getting-information-about-plugins"},"info"),": Info command is used to get suitable information about various plugins.\nSpecify the type of plugin as extractor, sink or processor.\nReturns information like, sample config, output and brief description of the plugin.")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"#linting-recipes"},"lint"),": used for validation of the recipes.\nHelps in avoiding any failure during running the meteor due to invalid recipe format.")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"#listing-all-the-plugins"},"list"),": used to state all the plugins of a certain type.")),(0,i.kt)("li",{parentName:"ul"},(0,i.kt)("p",{parentName:"li"},(0,i.kt)("a",{parentName:"p",href:"#running-recipes"},"run"),": the command is used for running the metadata extraction as per the instructions in the recipe.\nCan be used to run a single recipe, a directory of recipes or all the recipes in the current directory."))),(0,i.kt)("h2",{id:"listing-all-the-plugins"},"Listing all the plugins"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"# list all available extractors\n$ meteor list extractors\n\n# list all extractors with alias 'e'\n$ meteor list e\n\n# list available sinks\n$ meteor list sinks\n\n# list all sinks with alias 's'\n$ meteor list s\n\n# list all available processors\n$ meteor list processors\n\n# list all processors with alias 'p'\n$ meteor list p\n")),(0,i.kt)("h2",{id:"getting-information-about-plugins"},"Getting Information about plugins"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"# used to get info about different kinds of plugins\n$ meteor info [plugin-type] <plugin-name>\n\n# plugin-type can be sink, extractor or processor\n$ meteor info sink console\n$ meteor info processor enrich\n$ meteor info extractor postgres\n")),(0,i.kt)("h2",{id:"generating-sample-recipes"},"Generating Sample recipe","(","s",")"),(0,i.kt)("p",null,"Since recipe is the main resource of Meteor, we first need to create it before anything else.\nYou can create a sample recipe using the gen command."),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"# generate a sample recipe\n# generate a recipe with a bigquery extractor and a console sink\n$ meteor gen recipe sample -e bigquery -s console\n\n# generate recipe with multiple sinks\n$ meteor gen recipe sample -e bigquery -s compass,kafka\n\n# extractor(-e) as postgres, multiple sinks(-s) and enrich processor(-p)\n# save the generated recipe to a recipe.yaml\nmeteor gen recipe sample -e postgres -s compass,kafka -p enrich > recipe.yaml\n")),(0,i.kt)("h2",{id:"linting-recipes"},"Linting recipes"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"# validate specified recipes.\n$ meteor lint recipe.yml\n\n# lint all recipes in the specified directory\n$ meteor lint _recipes/\n\n# lint all recipes in the current directory\n$ meteor lint .\n")),(0,i.kt)("h2",{id:"running-recipes"},"Running recipes"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"# run meteor for specified recipes\n$ meteor run recipe.yml\n\n# run all recipes in the specified directory\n$ meteor run _recipes/\n\n# run all recipes in the current directory\n$ meteor run .\n")),(0,i.kt)("h2",{id:"get-help-on-commands-when-stuck"},"get help on commands when stuck"),(0,i.kt)("pre",null,(0,i.kt)("code",{parentName:"pre",className:"language-bash"},"# check for meteor help\n$ meteor --help\n")))}m.isMDXComponent=!0}}]);