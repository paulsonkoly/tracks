import { nodeResolve } from '@rollup/plugin-node-resolve';

// rollup.config.mjs
export default {
	input: 'tracks.js',
	output: {
		file: '../static/js/tracks.js',
		format: 'cjs'
	},
    plugins: [
		nodeResolve()
	]
};
