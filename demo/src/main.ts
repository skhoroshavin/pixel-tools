import {Game, Scale, Scene} from 'phaser';
const { ScaleModes } = Scale;

const textureAtlas = 'sprites'

export class GameScene extends Scene
{
    preload() {
        this.load.setPath('packed_assets')
        this.load.atlas(textureAtlas, textureAtlas+'.png', textureAtlas+'.atlas')
    }

    create() {
        const hero = this.createCharacter('hero')
        hero.anims.play('walk_front')

        const camera = this.cameras.main
        camera.setZoom(2)
        camera.startFollow(hero)
    }

    private createCharacter(name: string) {
        const res = this.add.sprite(0, 0, textureAtlas, `${name}_idle_front`)

        const idleAnim = (dir: string) =>
            res.anims.create({key: `idle_${dir}`, frames: [{key: textureAtlas, frame: `${name}_idle_${dir}`}]})
        idleAnim('front')
        idleAnim('back')
        idleAnim('left')
        idleAnim('right')

        const walkAnim = (dir: string) =>
            res.anims.create({
                key: `walk_${dir}`,
                frames: this.anims.generateFrameNames(textureAtlas, {
                    prefix: `${name}_walk_${dir}`,
                    zeroPad: 2,
                    end: 5
                }),
                frameRate: 8,
                repeat: -1
            })
        walkAnim('front')
        walkAnim('back')
        walkAnim('left')
        walkAnim('right')

        return res
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
