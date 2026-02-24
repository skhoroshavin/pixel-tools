import {Game, Scale, Scene} from 'phaser';
const { ScaleModes } = Scale;

export class GameScene extends Scene
{
    preload() {
        this.load.setPath('packed_assets')
    }

    create() {

    }
}

document.addEventListener('DOMContentLoaded', () => {
    new Game({
        parent: 'game-container',
        autoFocus: true,
        pixelArt: true,
        scale: {
            zoom: 3,
            mode: ScaleModes.RESIZE,
            autoCenter: Scale.CENTER_BOTH,
        },
        scene: GameScene,
    })
})
