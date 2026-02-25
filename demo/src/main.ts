import {Game, Scale, Scene} from 'phaser';
import Line = Phaser.Geom.Line;
const { ScaleModes } = Scale;

const textureAtlas = 'sprites'
const textureTilemap = 'tileset'

export class GameScene extends Scene
{
    preload() {
        this.load.setPath('packed_assets')
        this.load.atlas(textureAtlas, textureAtlas+'.png', textureAtlas+'.atlas')
        this.load.atlas(textureTilemap, textureTilemap+'.png', textureTilemap+'.atlas')
        this.load.tilemapTiledJSON('entrance', 'entrance.tmj')
        this.load.tilemapTiledJSON('cave', 'cave.tmj')
    }

    create() {
        const map = this.make.tilemap({key: 'entrance'})
        const tileset = map.addTilesetImage(textureTilemap)!
        map.layers.forEach((_, i) => map.createLayer(i, tileset))
        map.objects.forEach(l => l.objects.forEach(o => {
            if (o.gid === undefined) return
            const obj = this.add.sprite(o.x ?? 0, o.y ?? 0, textureTilemap, `${o.gid}`)
            obj.flipX = o.flippedHorizontal ?? false
            obj.flipY = o.flippedVertical ?? false
            obj.setOrigin(0, 1)
        }))

        const viking = this.createCharacter('viking')
        this.runSequence({
            sprite: viking,
            startAnim: 'walk_right',
            startX: map.widthInPixels/8,
            startY: 3*map.heightInPixels/4,
            endAnim: 'walk_left',
            endX: map.widthInPixels/3,
            endY: 3*map.heightInPixels/4,
        })

        const cat = this.createCharacter('cat')
        this.runSequence({
            sprite: cat,
            startAnim: 'walk_back',
            startX: 4*map.widthInPixels/7,
            startY: 5*map.heightInPixels/7,
            endAnim: 'walk_front',
            endX: 4*map.widthInPixels/7,
            endY: 4*map.heightInPixels/7,
        })

        const camera = this.cameras.main
        camera.setZoom(2)
        camera.setBounds(0, 0, map.widthInPixels, map.heightInPixels)
    }

    private createCharacter(name: string) {
        const res = this.add.sprite(0, 0, textureAtlas, `${name}_idle_front`)
        res.name = name

        const availableAnims = res.texture.getFrameNames()
            .filter(f => f.startsWith(name))
            .map(f => f.replace(name + '_', '').replace(/\d\d/,''))
            .sort()
            .filter((v, i, a) => i > 0 ? a[i-1] != v : true)

        availableAnims.forEach(anim => this.createAnimation(res, anim))
        return res
    }

    private createAnimation(sprite: Phaser.GameObjects.Sprite, name: string) {
        const prefix = `${sprite.name}_${name}`
        const frames = sprite.texture.getFrameNames().filter(f => f.startsWith(prefix)).sort()

        sprite.anims.create({
            key: name,
            frames: frames.map(f => ({key: textureAtlas, frame: f})),
            frameRate: 8,
            repeat: -1
        })
    }

    private runSequence(seq: {
        sprite: Phaser.GameObjects.Sprite,
        startAnim: string,
        startX: number,
        startY: number,
        endAnim: string,
        endX: number,
        endY: number,
    }) {
        seq.sprite.setPosition(seq.startX, seq.startY)
        const len = Phaser.Geom.Line.Length(new Line(seq.startX, seq.startY, seq.endX, seq.endY))
        const cfg = {targets: seq.sprite, duration: len * 12}
        const startTween = () => {
            seq.sprite.play(seq.startAnim)
            this.tweens.add({
                ...cfg,
                x: seq.endX,
                y: seq.endY,
                onComplete: endTween
            })
        }
        const endTween = () => {
            seq.sprite.play(seq.endAnim)
            this.tweens.add({
                ...cfg,
                x: seq.startX,
                y: seq.startY,
                onComplete: startTween
            })
        }
        startTween()
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new Game({
        parent: 'game-container',
        autoFocus: true,
        pixelArt: true,
        scale: {
            mode: ScaleModes.RESIZE,
            autoCenter: Scale.CENTER_BOTH,
        },
        scene: GameScene,
    })
})
