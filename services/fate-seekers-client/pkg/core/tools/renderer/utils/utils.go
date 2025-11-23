package utils

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/tile"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/google/uuid"
	"github.com/lafriks/go-tiled"
)

// PerformLoadMap loads map to different renderer levels calling callback when it ends.
func PerformLoadMap(tilemap *tiled.Map, callback func(spawnables []dto.Position)) {
	go func() {
		var spawnables []dto.Position

		for _, layer := range tilemap.Layers {
			layerTiles, spawnableTiles := loader.GetMapLayerTiles(
				layer, tilemap.Height, tilemap.Width, tilemap.TileHeight, tilemap.TileWidth)

			spawnables = append(spawnables, spawnableTiles...)

			layerTiles.Reverse(func(key float64, tiles []*dto.ProcessedTile) bool {
				for _, value := range tiles {
					name := uuid.New().String()

					switch layer.Name {
					case loader.FirstMapThirdLayer:
						if !renderer.GetInstance().TertiaryTileObjectExists(name) {
							renderer.GetInstance().AddTertiaryTileObject(
								name, tile.NewTile(value))
						}
					case loader.FirstMapSecondLayer:
						if !renderer.GetInstance().SecondaryTileObjectExists(name) {
							renderer.GetInstance().AddSecondaryTileObject(
								name, tile.NewTile(value))
						}
					}
				}

				return true
			})
		}

		callback(spawnables)
	}()
}
