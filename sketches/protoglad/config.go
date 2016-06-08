package main

import "fmt"

import "project/trigger/github.com/ehedgehog/griffin/project/trigger/rdf"

type Config struct {
	ShortnameToTerm map[string]rdf.Term
	Endpoints       []Endpoint
}

type Endpoint struct {
	UriTemplate string
	Term        rdf.Term
}

func (c Config) TermFor(shortname string) rdf.Term {
	return c.ShortnameToTerm[shortname]
}

func ParseConfig(g rdf.Graph) Config {
	apis := SubjectsWithProperty(g, rdf.Type, api_API)
	fmt.Printf("There are %v APIs in this config\n", len(apis))
	//
	shortnameToTerm := map[string]rdf.Term{}
	for _, P := range SubjectsWithProperty(g, API_label, rdf.ANY) {
		for _, L := range ObjectsOf(g, P, API_label) {
			shortnameToTerm[L.Spelling()] = P
		}
	}
	//
	endpoints := []Endpoint{}
	for _, api := range apis {
		for _, endpoint := range ObjectsOf(g, api, API_endpoint) {
			templates := ObjectsOf(g, endpoint, API_uriTemplate)
			if len(templates) > 1 {
				panic("multiple templates for endpoint")
			}
			endpoints = append(endpoints, Endpoint{templates[0].Spelling(), endpoint})
		}
	}
	//
	c := Config{ShortnameToTerm: shortnameToTerm, Endpoints: endpoints}
	return c
}

/*

API Spec:

	specificationURI = specification.getURI();
    	defaultPageSize = RDFUtils.getIntValue( specification, API.defaultPageSize, QueryParameter.DEFAULT_PAGE_SIZE );
		maxPageSize = RDFUtils.getIntValue( specification, API.maxPageSize, QueryParameter.MAX_PAGE_SIZE );
        describeThreshold = RDFUtils.getIntValue( specification, EXTRAS.describeThreshold, DEFAULT_DESCRIBE_THRESHOLD );
		prefixes = ExtractPrefixMapping.from(specification);
        sns = loadShortnames(specification, loader);
        dataSource = GetDataSource.sourceFromSpec( fm, specification, am );
        describeSources = extractDescribeSources( fm, am, specification, dataSource );
        primaryTopic = getStringValue(specification, FOAF.primaryTopic, null);
        defaultLanguage = getStringValue(specification, API.lang, null);
        base = getStringValue( specification, API.base, null );
        bindings.putAll( VariableExtractor.findAndBindVariables(specification) );
        factoryTable = RendererFactoriesSpec.createFactoryTable( specification );
        hasParameterBasedContentNegotiation = specification.hasProperty( API.contentNegotiation, API.parameterBased );
        extractEndpointSpecifications( specification );

API Endpoint Spec:

    	this.apiSpec = apiSpec;
    	wantsContext = endpoint.hasLiteral( EXTRAS.wantsContext, true );
    	bindings.putAll( apiSpec.bindings );
        bindings.putAll( VariableExtractor.findAndBindVariables( bindings, endpoint ) );
        defaultLanguage = getStringValue(endpoint, API.lang, apiSpec.getDefaultLanguage());
    	defaultPageSize = getIntValue( endpoint, API.defaultPageSize, apiSpec.defaultPageSize );
		maxPageSize = getIntValue( endpoint, API.maxPageSize, apiSpec.maxPageSize );
		cachePolicyName = getStringValue( endpoint, EXTRAS.cachePolicyName, "default" );
		parentApi = parent;
        name = endpoint.getLocalName();
        itemTemplate = getStringValue( endpoint, API.itemTemplate, null );
        uriTemplate = createURITemplate( endpoint );
        endpointResource = endpoint;
        describeThreshold = getIntValue( endpoint, EXTRAS.describeThreshold, apiSpec.describeThreshold );
        instantiateBaseQuery( endpoint );
        views = extractViews( endpoint );
        factoryTable = Render
*/
