"use strict";(self.webpackChunkmeteor=self.webpackChunkmeteor||[]).push([[53],{1109:e=>{e.exports=JSON.parse('{"pluginId":"default","version":"current","label":"Next","banner":null,"badge":false,"className":"docs-version-current","isLast":true,"docsSidebars":{"docsSidebar":[{"type":"link","label":"Introduction","href":"/meteor/docs/introduction","docId":"introduction"},{"type":"category","label":"Guides","items":[{"type":"link","label":"Introduction","href":"/meteor/docs/guides/introduction","docId":"guides/introduction"},{"type":"link","label":"Installation","href":"/meteor/docs/guides/installation","docId":"guides/installation"},{"type":"link","label":"Plugins","href":"/meteor/docs/guides/list_Plugins","docId":"guides/list_Plugins"},{"type":"link","label":"Recipes - Creation and linting","href":"/meteor/docs/guides/manage_recipes","docId":"guides/manage_recipes"},{"type":"link","label":"Running Meteor","href":"/meteor/docs/guides/run_recipes","docId":"guides/run_recipes"},{"type":"link","label":"Deployment","href":"/meteor/docs/guides/deployment","docId":"guides/deployment"}],"collapsed":false,"collapsible":true},{"type":"category","label":"Concepts","items":[{"type":"link","label":"Concepts","href":"/meteor/docs/concepts/overview","docId":"concepts/overview"},{"type":"link","label":"Recipe","href":"/meteor/docs/concepts/recipe","docId":"concepts/recipe"},{"type":"link","label":"Source","href":"/meteor/docs/concepts/source","docId":"concepts/source"},{"type":"link","label":"Processor","href":"/meteor/docs/concepts/processor","docId":"concepts/processor"},{"type":"link","label":"Sink","href":"/meteor/docs/concepts/sink","docId":"concepts/sink"}],"collapsed":false,"collapsible":true},{"type":"category","label":"Reference","items":[{"type":"link","label":"Commands","href":"/meteor/docs/reference/commands","docId":"reference/commands"},{"type":"link","label":"Configuration","href":"/meteor/docs/reference/configuration","docId":"reference/configuration"},{"type":"link","label":"Meteor Metadata Model","href":"/meteor/docs/reference/metadata_models","docId":"reference/metadata_models"},{"type":"link","label":"Extractors","href":"/meteor/docs/reference/extractors","docId":"reference/extractors"},{"type":"link","label":"Processors","href":"/meteor/docs/reference/processors","docId":"reference/processors"},{"type":"link","label":"Sinks","href":"/meteor/docs/reference/sinks","docId":"reference/sinks"}],"collapsed":false,"collapsible":true},{"type":"category","label":"Contribute","items":[{"type":"link","label":"Guide","href":"/meteor/docs/contribute/guide","docId":"contribute/guide"},{"type":"link","label":"Contribution Process","href":"/meteor/docs/contribute/contributing","docId":"contribute/contributing"}],"collapsed":false,"collapsible":true}]},"docs":{"concepts/overview":{"id":"concepts/overview","title":"Concepts","description":"A bit confused about various terms mentioned, and their usage?? Navigate through these concepts, this will help you writing a recipe for metadata extraction job:","sidebar":"docsSidebar"},"concepts/processor":{"id":"concepts/processor","title":"Processor","description":"A recipe can have none or many processors registered, depending upon the way the user wants metadata to be processed. A processor is basically a function that:","sidebar":"docsSidebar"},"concepts/recipe":{"id":"concepts/recipe","title":"Recipe","description":"A recipe is a set of instructions and configurations defined by the user, and in Meteor they are used to define how a particular job will be performed. It should contain instructions about the source from which the metadata will be fetched, information about metadata processors and the destination is to be defined as sinks of metadata.","sidebar":"docsSidebar"},"concepts/sink":{"id":"concepts/sink","title":"Sink","description":"sinks are used to define the medium of consuming the metadata being extracted. You need to specify at least one sink or can specify multiple sinks in a recipe, this will prevent you from having to create duplicate recipes for the same job. The given examples show you its correct usage if your sink is http and kafka.","sidebar":"docsSidebar"},"concepts/source":{"id":"concepts/source","title":"Source","description":"When the source field is defined, Meteor will extract data from a metadata","sidebar":"docsSidebar"},"contribute/contributing":{"id":"contribute/contributing","title":"Contribution Process","description":"BECOME A COMMITOR & CONTRIBUTE","sidebar":"docsSidebar"},"contribute/guide":{"id":"contribute/guide","title":"Guide","description":"Adding a new Extractor","sidebar":"docsSidebar"},"example/README":{"id":"example/README","title":"Example","description":"Running recipe with dynamic variable"},"guides/deployment":{"id":"guides/deployment","title":"Deployment","description":"After we are done with running and verifying that the recipes works with the data-source and sink you have mentioned.","sidebar":"docsSidebar"},"guides/installation":{"id":"guides/installation","title":"Installation","description":"Meteor can be installed currently by one of the following ways:","sidebar":"docsSidebar"},"guides/introduction":{"id":"guides/introduction","title":"Introduction","description":"The tour introduces you to meteor, the metadata orchestrator.","sidebar":"docsSidebar"},"guides/list_Plugins":{"id":"guides/list_Plugins","title":"Plugins","description":"Before getting started we expect you went through the prerequisites.","sidebar":"docsSidebar"},"guides/manage_recipes":{"id":"guides/manage_recipes","title":"Recipes - Creation and linting","description":"A recipe is a set of instructions and configurations defined by the user, and in Meteor they are used to define how a particular job will be performed.","sidebar":"docsSidebar"},"guides/run_recipes":{"id":"guides/run_recipes","title":"Running Meteor","description":"After we are done with creating sample recipe or a folder of sample recipes.","sidebar":"docsSidebar"},"introduction":{"id":"introduction","title":"Introduction","description":"Meteor is a plugin driven agent for collecting metadata. Meteor has plugins to source metadata from a variety of data stores, services and message queues. It also has sink plugins to send metadata to variety of third party APIs and catalog services.","sidebar":"docsSidebar"},"reference/commands":{"id":"reference/commands","title":"Commands","description":"Meteor currently supports the following commands and these can be utilised after the installation:","sidebar":"docsSidebar"},"reference/configuration":{"id":"reference/configuration","title":"Configuration","description":"This page contains references for all the application configurations for Meteor.","sidebar":"docsSidebar"},"reference/extractors":{"id":"reference/extractors","title":"Extractors","description":"Meteor currently supports metadata extraction on these data sources. To perform","sidebar":"docsSidebar"},"reference/metadata_models":{"id":"reference/metadata_models","title":"Meteor Metadata Model","description":"We have a set of defined metadata models which define the structure of metadata","sidebar":"docsSidebar"},"reference/processors":{"id":"reference/processors","title":"Processors","description":"Enrich","sidebar":"docsSidebar"},"reference/sinks":{"id":"reference/sinks","title":"Sinks","description":"Console","sidebar":"docsSidebar"}}}')}}]);