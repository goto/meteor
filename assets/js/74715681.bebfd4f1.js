"use strict";(self.webpackChunkmeteor=self.webpackChunkmeteor||[]).push([[21],{3905:function(t,e,a){a.d(e,{Zo:function(){return d},kt:function(){return c}});var r=a(7294);function n(t,e,a){return e in t?Object.defineProperty(t,e,{value:a,enumerable:!0,configurable:!0,writable:!0}):t[e]=a,t}function l(t,e){var a=Object.keys(t);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(t);e&&(r=r.filter((function(e){return Object.getOwnPropertyDescriptor(t,e).enumerable}))),a.push.apply(a,r)}return a}function i(t){for(var e=1;e<arguments.length;e++){var a=null!=arguments[e]?arguments[e]:{};e%2?l(Object(a),!0).forEach((function(e){n(t,e,a[e])})):Object.getOwnPropertyDescriptors?Object.defineProperties(t,Object.getOwnPropertyDescriptors(a)):l(Object(a)).forEach((function(e){Object.defineProperty(t,e,Object.getOwnPropertyDescriptor(a,e))}))}return t}function p(t,e){if(null==t)return{};var a,r,n=function(t,e){if(null==t)return{};var a,r,n={},l=Object.keys(t);for(r=0;r<l.length;r++)a=l[r],e.indexOf(a)>=0||(n[a]=t[a]);return n}(t,e);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(t);for(r=0;r<l.length;r++)a=l[r],e.indexOf(a)>=0||Object.prototype.propertyIsEnumerable.call(t,a)&&(n[a]=t[a])}return n}var m=r.createContext({}),o=function(t){var e=r.useContext(m),a=e;return t&&(a="function"==typeof t?t(e):i(i({},e),t)),a},d=function(t){var e=o(t.components);return r.createElement(m.Provider,{value:e},t.children)},f={inlineCode:"code",wrapper:function(t){var e=t.children;return r.createElement(r.Fragment,{},e)}},k=r.forwardRef((function(t,e){var a=t.components,n=t.mdxType,l=t.originalType,m=t.parentName,d=p(t,["components","mdxType","originalType","parentName"]),k=o(a),c=n,g=k["".concat(m,".").concat(c)]||k[c]||f[c]||l;return a?r.createElement(g,i(i({ref:e},d),{},{components:a})):r.createElement(g,i({ref:e},d))}));function c(t,e){var a=arguments,n=e&&e.mdxType;if("string"==typeof t||n){var l=a.length,i=new Array(l);i[0]=k;var p={};for(var m in e)hasOwnProperty.call(e,m)&&(p[m]=e[m]);p.originalType=t,p.mdxType="string"==typeof t?t:n,i[1]=p;for(var o=2;o<l;o++)i[o]=a[o];return r.createElement.apply(null,i)}return r.createElement.apply(null,a)}k.displayName="MDXCreateElement"},8971:function(t,e,a){a.r(e),a.d(e,{frontMatter:function(){return p},contentTitle:function(){return m},metadata:function(){return o},toc:function(){return d},default:function(){return k}});var r=a(7462),n=a(3366),l=(a(7294),a(3905)),i=["components"],p={},m="Extractors",o={unversionedId:"reference/extractors",id:"reference/extractors",isDocsHomePage:!1,title:"Extractors",description:"Meteor currently support metadata extraction on these data sources. To perform extraction on any of these you need to create a recipe file with instructions as mentioned here. In the sample-recipe.yaml add source information such as type from the table below and config for that particular extractor can be found by visiting the link in type field.",source:"@site/docs/reference/extractors.md",sourceDirName:"reference",slug:"/reference/extractors",permalink:"/meteor/docs/reference/extractors",editUrl:"https://github.com/odpf/meteor/edit/master/docs/docs/reference/extractors.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Meteor Metadata Model",permalink:"/meteor/docs/reference/metadata_models"},next:{title:"Processors",permalink:"/meteor/docs/reference/processors"}},d=[{value:"Extractors Feature Matrix",id:"extractors-feature-matrix",children:[{value:"Table",id:"table",children:[]},{value:"Dashboard",id:"dashboard",children:[]},{value:"Topic",id:"topic",children:[]},{value:"User",id:"user",children:[]},{value:"Bucket",id:"bucket",children:[]},{value:"Job",id:"job",children:[]}]}],f={toc:d};function k(t){var e=t.components,a=(0,n.Z)(t,i);return(0,l.kt)("wrapper",(0,r.Z)({},f,a,{components:e,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"extractors"},"Extractors"),(0,l.kt)("p",null,"Meteor currently support metadata extraction on these data sources. To perform extraction on any of these you need to create a recipe file with instructions as mentioned ",(0,l.kt)("a",{parentName:"p",href:"/meteor/docs/concepts/recipe"},"here"),". In the ",(0,l.kt)("inlineCode",{parentName:"p"},"sample-recipe.yaml")," add ",(0,l.kt)("inlineCode",{parentName:"p"},"source")," information such as ",(0,l.kt)("inlineCode",{parentName:"p"},"type")," from the table below and ",(0,l.kt)("inlineCode",{parentName:"p"},"config")," for that particular extractor can be found by visiting the link in ",(0,l.kt)("inlineCode",{parentName:"p"},"type")," field."),(0,l.kt)("h2",{id:"extractors-feature-matrix"},"Extractors Feature Matrix"),(0,l.kt)("h3",{id:"table"},"Table"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:"left"},"Type"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Attributes"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Profile"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Schema"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Lineage"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Ownership"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Custom"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/clickhouse/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"clickhouse"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/couchdb/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"couchdb"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/mongodb/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"mongodb"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/mssql/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"mssql"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/mysql/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"mysql"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/postgres/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"postgres"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/cassandra/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"cassandra"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")))),(0,l.kt)("h3",{id:"dashboard"},"Dashboard"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:"left"},"Type"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Url"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Chart"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Lineage"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Tags"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Custom"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/cb12c3ecf8904cf3f4ce365ca8981ccd132f35d0/plugins/extractors/grafana/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"grafana"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/cb12c3ecf8904cf3f4ce365ca8981ccd132f35d0/plugins/extractors/metabase/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"metabase"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/cb12c3ecf8904cf3f4ce365ca8981ccd132f35d0/plugins/extractors/superset/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"superset"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")),(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/cb12c3ecf8904cf3f4ce365ca8981ccd132f35d0/plugins/extractors/tableau/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"tableau"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")))),(0,l.kt)("h3",{id:"topic"},"Topic"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:"left"},"Type"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Profile"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Schema"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Ownership"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Lineage"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Tags"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Custom"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/kafka/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"kafka"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717")))),(0,l.kt)("h3",{id:"user"},"User"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:"left"},"Type"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Email"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Username"),(0,l.kt)("th",{parentName:"tr",align:"left"},"FullName"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Title"),(0,l.kt)("th",{parentName:"tr",align:"left"},"IsActive"),(0,l.kt)("th",{parentName:"tr",align:"left"},"ManagerEmail"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Profiles"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Memberships"),(0,l.kt)("th",{parentName:"tr",align:"left"},"facets"),(0,l.kt)("th",{parentName:"tr",align:"left"},"common"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/github/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"github"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2610"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2610"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2610"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2610"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2610"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2610")))),(0,l.kt)("h3",{id:"bucket"},"Bucket"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:"left"},"type"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Location"),(0,l.kt)("th",{parentName:"tr",align:"left"},"StorageType"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Blobs"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Ownership"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Tags"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Custom"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Timestamps"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/gcs/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"gcs"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2717"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705")))),(0,l.kt)("h3",{id:"job"},"Job"),(0,l.kt)("table",null,(0,l.kt)("thead",{parentName:"table"},(0,l.kt)("tr",{parentName:"thead"},(0,l.kt)("th",{parentName:"tr",align:"left"},"Type"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Ownership"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Upstreams"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Downstreams"),(0,l.kt)("th",{parentName:"tr",align:"left"},"Custom"))),(0,l.kt)("tbody",{parentName:"table"},(0,l.kt)("tr",{parentName:"tbody"},(0,l.kt)("td",{parentName:"tr",align:"left"},(0,l.kt)("a",{parentName:"td",href:"https://github.com/odpf/meteor/tree/main/plugins/extractors/optimus/README.md"},(0,l.kt)("inlineCode",{parentName:"a"},"optimus"))),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705"),(0,l.kt)("td",{parentName:"tr",align:"left"},"\u2705")))))}k.isMDXComponent=!0}}]);