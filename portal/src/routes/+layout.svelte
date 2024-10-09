<script lang="ts">
	import { DarkMode } from 'flowbite-svelte';
	import { Navbar, NavBrand, NavLi, NavUl, NavHamburger } from 'flowbite-svelte';
	import { page } from '$app/stores';
	import '../app.css';

	function getPageTitle(path: string) {
		if (path === '') {
			return 'Kycelis';
		}

		if (path === '/') {
			return 'Dashboard';
		}

		let title: string = path.replace(/-/g, ' ');
		title = title.charAt(1).toUpperCase() + title.slice(2);

		return title;
	}

	let pageTitle: string;
	$: pageTitle = getPageTitle($page.url.pathname);
</script>

<svelte:head>
	<title>{pageTitle}</title>
</svelte:head>

<Navbar class="bg-slate-200 dark:bg-slate-900">
	<NavBrand href="/">
		<img src="/logo.png" class="me-3 h-6 sm:h-9" alt="Flowbite Logo" />
		<span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">Kycelis</span>
	</NavBrand>
	<NavHamburger />
	<NavUl>
		<NavLi href="/">Dashboard</NavLi>
		<NavLi href="/topology">Topology</NavLi>
		<NavLi href="/inventory">Inventory</NavLi>
	</NavUl>
	<DarkMode />
</Navbar>

<slot></slot>
