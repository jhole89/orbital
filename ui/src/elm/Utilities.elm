module Utilities exposing (..)

import ECharts exposing (ChartOptions, GraphEdgeItemOption, GraphNodeItemOption)
import EntityHelpers exposing (Entity, EntityListResponse)
import List.Extra


entityListToChartOpts : EntityListResponse -> ChartOptions
entityListToChartOpts entityList =
    let
        categories = List.sort (List.Extra.unique (List.map .context entityList))
    in
        { title = Nothing
        , legend = [{ data = categories }]
        , tooltip =
            Just { show = False
            }
        , series = [
            { animation = Just True
            , categories = List.map (\c -> { name = c }) categories
            , data = List.map (entityToGraphNodeItemOption categories) entityList
            , draggable = Just True
            , emphasis = Nothing
            , force =
                Just
                    { edgeLength = Just 100
                    , friction = Just 0.2
                    , gravity = Just 0.2
                    , layoutAnimation = Just True
                    , repulsion = Just 50
                    }
            , label =
                Just
                    { position = Just "right"
                    , formatter = Just "{b}"
                    , show = Nothing
                    }
            , layout = "force"
            , lineStyle = Nothing
            , links = List.concatMap entityToGraphEdgeItemOption entityList
            , name = Nothing
            , roam = Just True
            , type_ = "graph"
            }]
        }


entityToGraphNodeItemOption : List String -> Entity -> GraphNodeItemOption
entityToGraphNodeItemOption categories entity =
    { name = Just entity.name
    , id = Just entity.id
    , category = List.Extra.elemIndex entity.context categories
    , value = 1
    , x = Nothing
    , y = Nothing
    , symbolSize = Just 20
    , label = Nothing
    }

entityToGraphEdgeItemOption: Entity -> List GraphEdgeItemOption
entityToGraphEdgeItemOption e =
    List.map (\c -> { source = e.name, target = c}) e.connections

