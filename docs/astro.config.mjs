// @ts-check
import {defineConfig} from 'astro/config';
import starlight from '@astrojs/starlight';
import catppuccin from "starlight-theme-catppuccin";

// https://astro.build/config
export default defineConfig({
	site: "https://quentinrob.github.io",
	base: "pkgs-manager",
	integrations: [
		starlight({
			title: 'PkgsManager',
			social: [{icon: 'github', label: 'GitHub', href: 'https://github.com/QuentinRob/PkgsManager'}],
			sidebar: [
				{
					label: 'Guides',
					items: [
						// Each item here is one entry in the navigation menu.
						{label: 'Example Guide', slug: 'guides/example'},
					],
				},
				{
					label: 'Reference',
					autogenerate: {directory: 'reference'},
				},
			],
			plugins: [
				// @ts-ignore
				catppuccin()
			]
		}),
	],
});
