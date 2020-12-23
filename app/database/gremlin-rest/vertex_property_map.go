package gremlin_rest


type VertexPropertyMapContainer struct {
	Type  string         `json:"@type"`
	Value VertexProperty `json:"@value"`
}
type VertexPropertyMapListContainer struct {
	Type string                     `json:"@type"`
	Value []VertexPropertyMapContainer `json:"@value"`
}


// {
//  "@type": "g:List",
//  "@value": [
//    {
//      "@type": "g:Map",
//      "@value": [
//        "name",
//        {
//          "@type": "g:List",
//          "@value": [
//            {
//              "@type": "g:VertexProperty",
//              "@value": {
//                "id": {
//                  "@type": "g:Int64",
//                  "@value": 799
//                },
//                "value": "active",
//                "label": "name"
//              }
//            }
//          ]
//        },
//        "context",
//        {
//          "@type": "g:List",
//          "@value": [
//            {
//              "@type": "g:VertexProperty",
//              "@value": {
//                "id": {
//                  "@type": "g:Int64",
//                  "@value": 800
//                },
//                "value": "field",
//                "label": "context"
//              }
//            }
//          ]
//        }
//      ]
//    },
//