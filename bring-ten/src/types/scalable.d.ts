declare module 'scalable' {
	interface ScalableOptions {
		align?: 'left' | 'center' | 'right';
		verticalAlign?: 'top' | 'center' | 'bottom';
		containerHeight?: 'fixed' | 'auto';
		minWidth?: number;
		maxWidth?: number;
		minScale?: number;
		maxScale?: number;
	}

	export default class Scalable {
		constructor(container: HTMLElement, options?: ScalableOptions);
		destroy?(): void;
		update(): void;
	}
}
