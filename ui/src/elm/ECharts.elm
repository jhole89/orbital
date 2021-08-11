module ECharts exposing (..)

import Entity.ScalarCodecs exposing (Id)
import Json.Encode as Encode


type alias Data =
    { nodes : List GraphNodeItemOption
    , links : List GraphEdgeItemOption
    , categories : List CategoryItemOption
    }


type alias GraphEdgeItemOption =
    { source : String
    , target : String
    }


type alias GraphNodeItemOption =
    { name : Maybe String
    , id : Maybe Id
    , value : Float
    , x : Maybe Float
    , y : Maybe Float
    , category : Maybe Int
    , symbolSize : Maybe Float
    , label : Maybe LabelOption
    }


type alias CategoryItemOption =
    { name : String
    }


type alias ChartOptions =
    { title : Maybe TitleOption
    , tooltip : Maybe TooltipOption
    , legend : List LegendOption
    , series : List SeriesOption
    }


type alias TitleOption =
    { text : String
    , subtext : String
    , top : String
    , left : String
    }


type alias TooltipOption =
    { show : Bool
    }


type alias LegendOption =
    { data : List String
    }


type alias SeriesOption =
    { animation : Maybe Bool
    , categories : List CategoryItemOption
    , data : List GraphNodeItemOption
    , draggable : Maybe Bool
    , emphasis : Maybe EmphasisOption
    , force : Maybe ForceOption
    , label : Maybe LabelOption
    , layout : String
    , lineStyle : Maybe LineStyleOption
    , links : List GraphEdgeItemOption
    , name : Maybe String
    , roam : Maybe Bool
    , type_ : String
    }


type alias LabelOption =
    { position : Maybe String
    , formatter : Maybe String
    , show : Maybe Bool
    }


type alias LineStyleOption =
    { color : Maybe String
    , curveness : Maybe Float
    , width : Maybe Int
    }


type alias EmphasisOption =
    { focus : String
    , lineStyle : LineStyleOption
    }


type alias ForceOption =
    { edgeLength : Maybe Int
    , friction : Maybe Float
    , gravity : Maybe Float
    , layoutAnimation : Maybe Bool
    , repulsion : Maybe Int
    }


encodeChartOptions : ChartOptions -> Encode.Value
encodeChartOptions chartOptions =
    (Encode.object << List.filterMap identity)
        [ chartOptions.title |> Maybe.andThen (\title -> Just ( "title", titleOptionEncoder title ))
        , chartOptions.tooltip |> Maybe.andThen (\t -> Just ( "tooltip", tooltipOptionEncoder t ))
        , Just ( "legend", Encode.list legendOptionEncoder chartOptions.legend )
        , Just ( "series", Encode.list seriesOptionEncoder chartOptions.series )
        ]


titleOptionEncoder : TitleOption -> Encode.Value
titleOptionEncoder titleOption =
    Encode.object
        [ ( "text", Encode.string titleOption.text )
        , ( "subtext", Encode.string titleOption.subtext )
        , ( "top", Encode.string titleOption.top )
        , ( "left", Encode.string titleOption.left )
        ]


tooltipOptionEncoder : TooltipOption -> Encode.Value
tooltipOptionEncoder tooltipOption =
    Encode.object
        [ ( "show", Encode.bool tooltipOption.show )
        ]


legendOptionEncoder : LegendOption -> Encode.Value
legendOptionEncoder legendOption =
    Encode.object [ ( "data", Encode.list Encode.string legendOption.data ) ]


seriesOptionEncoder : SeriesOption -> Encode.Value
seriesOptionEncoder seriesOption =
    (Encode.object << List.filterMap identity)
        [ seriesOption.animation |> Maybe.andThen (\a -> Just ( "animation", Encode.bool a ))
        , Just ( "categories", Encode.list categoryItemOptionEncoder seriesOption.categories )
        , Just ( "data", Encode.list graphNodeItemOptionEncoder seriesOption.data )
        , seriesOption.draggable |> Maybe.andThen (\d -> Just ( "draggable", Encode.bool d ))
        , seriesOption.emphasis |> Maybe.andThen (\e -> Just ( "emphasis", emphasisOptionEncoder e ))
        , seriesOption.force |> Maybe.andThen (\f -> Just ( "force", forceOptionEncoder f ))
        , seriesOption.label |> Maybe.andThen (\l -> Just ( "label", labelOptionEncoder l ))
        , Just ( "layout", Encode.string seriesOption.layout )
        , seriesOption.lineStyle |> Maybe.andThen (\ls -> Just ( "lineStyle", lineStyleOptionEncoder ls ))
        , Just ( "links", Encode.list graphEdgeItemOptionEncoder seriesOption.links )
        , seriesOption.name |> Maybe.andThen (\n -> Just ( "name", Encode.string n ))
        , seriesOption.roam |> Maybe.andThen (\r -> Just ( "roam", Encode.bool r ))
        , Just ( "type", Encode.string seriesOption.type_ )
        ]


labelOptionEncoder : LabelOption -> Encode.Value
labelOptionEncoder labelOption =
    (Encode.object << List.filterMap identity)
        [ labelOption.position |> Maybe.andThen (\p -> Just ( "position", Encode.string p ))
        , labelOption.formatter |> Maybe.andThen (\f -> Just ( "formatter", Encode.string f ))
        , labelOption.show |> Maybe.andThen (\s -> Just ( "show", Encode.bool s ))
        ]


lineStyleOptionEncoder : LineStyleOption -> Encode.Value
lineStyleOptionEncoder lineStyleOption =
    (Encode.object << List.filterMap identity)
        [ lineStyleOption.color
            |> Maybe.andThen (\color -> Just ( "color", Encode.string color ))
        , lineStyleOption.curveness
            |> Maybe.andThen (\curveness -> Just ( "curveness", Encode.float curveness ))
        , lineStyleOption.width
            |> Maybe.andThen (\width -> Just ( "width", Encode.int width ))
        ]


emphasisOptionEncoder : EmphasisOption -> Encode.Value
emphasisOptionEncoder emphasisOption =
    Encode.object
        [ ( "focus", Encode.string emphasisOption.focus )
        , ( "lineStyle", lineStyleOptionEncoder emphasisOption.lineStyle )
        ]


forceOptionEncoder : ForceOption -> Encode.Value
forceOptionEncoder forceOption =
    (Encode.object << List.filterMap identity)
        [ forceOption.edgeLength |> Maybe.andThen (\el -> Just ( "edgeLength", Encode.int el ))
        , forceOption.friction |> Maybe.andThen (\f -> Just ( "friction", Encode.float f ))
        , forceOption.gravity |> Maybe.andThen (\g -> Just ( "gravity", Encode.float g ))
        , forceOption.layoutAnimation |> Maybe.andThen (\la -> Just ( "layoutAnimation", Encode.bool la ))
        , forceOption.repulsion |> Maybe.andThen (\r -> Just ( "repulsion", Encode.int r ))
        ]


graphNodeItemOptionEncoder : GraphNodeItemOption -> Encode.Value
graphNodeItemOptionEncoder gnio =
    (Encode.object << List.filterMap identity)
        [ gnio.name |> Maybe.andThen (\n -> Just ( "name", Encode.string n ))
        , Just ( "value", Encode.float gnio.value )
        , gnio.x |> Maybe.andThen (\x -> Just ( "x", Encode.float x ))
        , gnio.y |> Maybe.andThen (\y -> Just ( "y", Encode.float y ))
        , gnio.category |> Maybe.andThen (\c -> Just ( "category", Encode.int c ))
        , gnio.symbolSize |> Maybe.andThen (\s -> Just ( "symbolSize", Encode.float s ))
        , gnio.label |> Maybe.andThen (\l -> Just ( "label", labelOptionEncoder l ))
        ]


graphEdgeItemOptionEncoder : GraphEdgeItemOption -> Encode.Value
graphEdgeItemOptionEncoder graphEdgeItemOption =
    Encode.object
        [ ( "source", Encode.string graphEdgeItemOption.source )
        , ( "target", Encode.string graphEdgeItemOption.target )
        ]


categoryItemOptionEncoder : CategoryItemOption -> Encode.Value
categoryItemOptionEncoder categoryItemOption =
    Encode.object
        [ ( "name", Encode.string categoryItemOption.name ) ]
