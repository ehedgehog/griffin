package rdf

const RDFS_NS = `http://www.w3.org/2000/01/rdf-schema#`
const RDF_NS = `http://www.w3.org/1999/02/22-rdf-syntax-ns#`

const Resource = IRI(RDFS_NS + `Resource`)
const Class = IRI(RDFS_NS + `Class`)
const Literal = IRI(RDFS_NS + `Literal`)
const Datatype = IRI(RDFS_NS + `Datatype`)
const Container = IRI(RDFS_NS + `Container`)
const ContainerMembershipProperty = IRI(RDFS_NS + `ContainerMembershipProperty`)

const SeeAlso = IRI(RDFS_NS + `seeAlso`)
const IsDefinedBy = IRI(RDFS_NS + `isDefinedBy`)
const Member = IRI(RDFS_NS + `member`)
const Range = IRI(RDFS_NS + `range`)
const Domain = IRI(RDFS_NS + `domain`)
const SubClassOf = IRI(RDFS_NS + `subClassOf`)
const SubPropertyOf = IRI(RDFS_NS + `subPropertyOf`)
const Label = IRI(RDFS_NS + `label`)
const Comment = IRI(RDFS_NS + `comment`)

const XMLLiteral = IRI(RDF_NS + `XMLLiteral`)

const Property = IRI(RDF_NS + `Property`)
const Statement = IRI(RDF_NS + `Statement`)
const Bag = IRI(RDF_NS + `Bag`)
const Seq = IRI(RDF_NS + `Seq`)
const Alt  = IRI(RDF_NS + `Alt `)
const Type = IRI(RDF_NS + `Type`)
const First = IRI(RDF_NS + `First`)
const Rest = IRI(RDF_NS + `Rest`)
const Nil = IRI(RDF_NS + `nil`)
const Value = IRI(RDF_NS + `Value`)
const Subject = IRI(RDF_NS + `Subject`)
const Predicate = IRI(RDF_NS + `Predicate`)
const Object = IRI(RDF_NS + `Object`)
