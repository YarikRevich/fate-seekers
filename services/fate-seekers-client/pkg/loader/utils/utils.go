package utils

import (
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/collision"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/renderer/tile"
	selected "github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/selecter"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/core/tools/sounder"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/dto"
	"github.com/YarikRevich/fate-seekers/services/fate-seekers-client/pkg/loader"
	"github.com/google/uuid"
	"github.com/lafriks/go-tiled"
)

// PerformLoadMap loads map to different renderer levels calling callback when it ends.
func PerformLoadMap(tilemap *tiled.Map, callback func(spawnables, chestLocations, healthPackLocations []dto.Position)) {
	go func() {
		var (
			spawnables          []dto.Position
			chestLocations      []dto.Position
			healthPackLocations []dto.Position
		)

		for _, layer := range tilemap.Layers {
			layerTiles,
				spawnableTiles,
				chestTiles,
				healthPackTiles,
				collidableTiles,
				soundableTiles,
				selectableTiles := loader.GetMapLayerTiles(
				layer, tilemap.Height, tilemap.Width, tilemap.TileHeight, tilemap.TileWidth)

			spawnables = append(spawnables, spawnableTiles...)
			chestLocations = append(spawnables, chestTiles...)
			healthPackLocations = append(spawnables, healthPackTiles...)

			for _, soundableTile := range soundableTiles {
				sounder.GetInstance().AddSoundableTileObject(soundableTile)
			}

			for _, collidableTile := range collidableTiles {
				collision.GetInstance().AddCollidableTileObject(collidableTile)
			}

			for _, selectableTile := range selectableTiles {
				selected.GetInstance().AddSelectableTileObject(selectableTile)
			}

			layerTiles.Reverse(func(key float64, tiles []*dto.ProcessedTile) bool {
				for _, value := range tiles {
					name := uuid.NewString()

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

		callback(spawnables, chestLocations, healthPackLocations)
	}()
}
