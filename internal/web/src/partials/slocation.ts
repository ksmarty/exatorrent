import { readable, type Readable } from 'svelte/store';

/* shape of the store value (mirrors Location) */
export interface SLocation extends Location {}

/* readable store that always reflects window.location */
const store = readable<SLocation>(typeof window !== 'undefined' ? { ...window.location } : ({} as SLocation), (set) => {
  const sync = () => set({ ...window.location });
  sync(); // initial value
  addEventListener('popstate', sync);
  addEventListener('hashchange', sync);
  return () => {
    removeEventListener('popstate', sync);
    removeEventListener('hashchange', sync);
  };
});

/* navigation helpers that keep the store reactive */
const goto = (url: string = '', replace: boolean = false): void => {
  history[replace ? 'replaceState' : 'pushState']({}, '', url);
  dispatchEvent(new PopStateEvent('popstate'));
};

const pushState = (data: any, title: string, url?: string | null): void => {
  history.pushState(data, title, url ?? '');
  dispatchEvent(new PopStateEvent('popstate'));
};

const replaceState = (data: any, title: string, url?: string | null): void => {
  history.replaceState(data, title, url ?? '');
  dispatchEvent(new PopStateEvent('popstate'));
};

const reset = (): void => {
  dispatchEvent(new PopStateEvent('popstate'));
};

/* single object that carries both the store and the methods */
export const slocation = {
  subscribe: store.subscribe, // Readable contract
  goto,
  pushState,
  replaceState,
  reset
} as Readable<SLocation> & {
  goto: typeof goto;
  pushState: typeof pushState;
  replaceState: typeof replaceState;
  reset: typeof reset;
};

/* default export for `import slocation from 'slocation'` */
export default slocation;
