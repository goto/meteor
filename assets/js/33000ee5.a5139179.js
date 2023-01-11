"use strict";(self.webpackChunkmeteor=self.webpackChunkmeteor||[]).push([[362],{3905:(e,t,a)=>{a.d(t,{Zo:()=>c,kt:()=>u});var r=a(7294);function n(e,t,a){return t in e?Object.defineProperty(e,t,{value:a,enumerable:!0,configurable:!0,writable:!0}):e[t]=a,e}function o(e,t){var a=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),a.push.apply(a,r)}return a}function s(e){for(var t=1;t<arguments.length;t++){var a=null!=arguments[t]?arguments[t]:{};t%2?o(Object(a),!0).forEach((function(t){n(e,t,a[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(a)):o(Object(a)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(a,t))}))}return e}function i(e,t){if(null==e)return{};var a,r,n=function(e,t){if(null==e)return{};var a,r,n={},o=Object.keys(e);for(r=0;r<o.length;r++)a=o[r],t.indexOf(a)>=0||(n[a]=e[a]);return n}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)a=o[r],t.indexOf(a)>=0||Object.prototype.propertyIsEnumerable.call(e,a)&&(n[a]=e[a])}return n}var p=r.createContext({}),l=function(e){var t=r.useContext(p),a=t;return e&&(a="function"==typeof e?e(t):s(s({},t),e)),a},c=function(e){var t=l(e.components);return r.createElement(p.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},d=r.forwardRef((function(e,t){var a=e.components,n=e.mdxType,o=e.originalType,p=e.parentName,c=i(e,["components","mdxType","originalType","parentName"]),d=l(a),u=n,f=d["".concat(p,".").concat(u)]||d[u]||m[u]||o;return a?r.createElement(f,s(s({ref:t},c),{},{components:a})):r.createElement(f,s({ref:t},c))}));function u(e,t){var a=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var o=a.length,s=new Array(o);s[0]=d;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i.mdxType="string"==typeof e?e:n,s[1]=i;for(var l=2;l<o;l++)s[l]=a[l];return r.createElement.apply(null,s)}return r.createElement.apply(null,a)}d.displayName="MDXCreateElement"},2038:(e,t,a)=>{a.r(t),a.d(t,{assets:()=>p,contentTitle:()=>s,default:()=>m,frontMatter:()=>o,metadata:()=>i,toc:()=>l});var r=a(7462),n=(a(7294),a(3905));const o={},s="Meteor Metadata Model",i={unversionedId:"reference/metadata_models",id:"reference/metadata_models",title:"Meteor Metadata Model",description:"We have a set of defined metadata models which define the structure of metadata",source:"@site/docs/reference/metadata_models.md",sourceDirName:"reference",slug:"/reference/metadata_models",permalink:"/meteor/docs/reference/metadata_models",draft:!1,editUrl:"https://github.com/odpf/meteor/edit/master/docs/docs/reference/metadata_models.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Configuration",permalink:"/meteor/docs/reference/configuration"},next:{title:"Extractors",permalink:"/meteor/docs/reference/extractors"}},p={},l=[{value:"Usage",id:"usage",level:2}],c={toc:l};function m(e){let{components:t,...a}=e;return(0,n.kt)("wrapper",(0,r.Z)({},c,a,{components:t,mdxType:"MDXLayout"}),(0,n.kt)("h1",{id:"meteor-metadata-model"},"Meteor Metadata Model"),(0,n.kt)("p",null,"We have a set of defined metadata models which define the structure of metadata\nthat meteor will yield. To visit the metadata models being used by different\nextractors please visit ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors"},"here"),". We are currently using the\nfollowing metadata models:"),(0,n.kt)("ul",null,(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/bucket.proto"},"Bucket"),": Used for metadata being extracted from buckets.\nBuckets are the basic containers in google cloud services, or Amazon S3, etc\nthat are used fot data storage, and quite popular because of their features of\naccess management, aggregation of usage and services and ease of\nconfigurations. Currently, Meteor provides a metadata extractor for the\nbuckets mentioned ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#bucket"},"here"))),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/dashboard.proto"},"Dashboard"),": Dashboards are an essential part of data\nanalysis and are used to track, analyze and visualize. These Dashboard\nmetadata model includes some basic fields like ",(0,n.kt)("inlineCode",{parentName:"p"},"urn")," and ",(0,n.kt)("inlineCode",{parentName:"p"},"source"),", etc and a\nlist of ",(0,n.kt)("inlineCode",{parentName:"p"},"Chart"),". There are multiple dashboards that are essential for Data\nAnalysis such as metabase, grafana, tableau, etc. Please refer to the list of\n'Dashboard' extractors meteor currently\nsupports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#dashboard"},"here"),"."),(0,n.kt)("ul",{parentName:"li"},(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("a",{parentName:"li",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/dashboard.proto"},"Chart"),": Charts are included in all the Dashboard and are\nthe result of certain queries in a Dashboard. Information about them\nincludes the information of the query and few similar details."))),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/user.proto"},"User"),": This metadata model is used for defining the output of\nextraction on User accounts. Some of these sources can be GitHub, Workday,\nGoogle Suite, LDAP. Please refer to the list of 'User' extractors meteor\ncurrently supports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#user"},"here"),".")),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/table.proto"},"Table"),": This metadata model is being used by extractors based\naround databases, typically for the ones that store data in tabular format. It\ncontains various fields that include ",(0,n.kt)("inlineCode",{parentName:"p"},"schema")," of the table and other access\nrelated information. Please refer to the list of 'Table' extractors meteor\ncurrently supports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#table"},"here"),".")),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/job.proto"},"Job"),": A job can represent a scheduled or recurring task that\nperforms some transformation in the data engineering pipeline. Job is a\nmetadata model built for this purpose. Please refer to the list of 'Job'\nextractors meteor currently supports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#table"},"here"),".")),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/topic.proto"},"Topic"),": A topic represents a virtual group for logical group of\nmessages in message bus like kafka, pubsub, pulsar etc. Please refer to the\nlist of 'Topic' extractors meteor currently\nsupports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#topic"},"here"),".")),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/feature_table.proto"},"Machine Learning Feature Table"),": A Feature Table is a\ntable or view that represents a logical group of time-series feature data as\nit is found in a data source. Please refer to the list of 'Feature Table'\nextractors meteor currently\nsupports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#machine-learning-feature-table"},"here"),".")),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/application.proto"},"Application"),": An application represents a service that\ntypically communicates over well-defined APIs. Please refer to the list of '\nApplication' extractors meteor currently\nsupports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#application"},"here"),".")),(0,n.kt)("li",{parentName:"ul"},(0,n.kt)("p",{parentName:"li"},(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2/model.proto"},"Machine Learning Model"),": A Model represents a Data Science\nModel commonly used for Machine Learning(ML). Models are algorithms trained on\ndata to find patterns or make predictions. Models typically consume ML\nfeatures to generate a meaningful output. Please refer to the list of 'Model'\nextractors meteor currently\nsupports ",(0,n.kt)("a",{parentName:"p",href:"/meteor/docs/reference/extractors#machine-learning-model"},"here"),"."))),(0,n.kt)("p",null,(0,n.kt)("inlineCode",{parentName:"p"},"Proto")," has been used to define these metadata models. To check their\nimplementation please refer ",(0,n.kt)("a",{parentName:"p",href:"https://github.com/odpf/proton/tree/main/odpf/assets/v1beta2"},"here"),"."),(0,n.kt)("h2",{id:"usage"},"Usage"),(0,n.kt)("pre",null,(0,n.kt)("code",{parentName:"pre",className:"language-golang"},'import(\n    assetsv1beta1 "github.com/odpf/meteor/models/odpf/assets/v1beta1"\n    "github.com/odpf/meteor/models/odpf/assets/facets/v1beta1"\n)\n\nfunc main(){\n    // result is a var of data type of assetsv1beta1.Table one of our metadata model\n    result := &assetsv1beta1.Table{\n        // assigining value to metadata model\n        Urn:  fmt.Sprintf("%s.%s", dbName, tableName),\n        Name: tableName,\n    }\n\n    // using column facet to add metadata info of schema\n\n    var columns []*facetsv1beta1.Column\n    columns = append(columns, &facetsv1beta1.Column{\n            Name:       "column_name",\n            DataType:   "varchar",\n            IsNullable: true,\n            Length:     256,\n        })\n    result.Schema = &facetsv1beta1.Columns{\n        Columns: columns,\n    }\n}\n')))}m.isMDXComponent=!0}}]);