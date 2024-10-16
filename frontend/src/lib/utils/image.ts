import Placeholder from "$assets/placeholders/16x9.png";

// For some reason images are rendered twice (or maybe more?)
// nullifying the 'onerror' on the first render...
// So take that into account, try thrice and if not, safely stop
export const imgPlaceholder: string = `
  this._count = this._count ? this._count + 1 : 1;
  this._count <= 3 ? this.src = '${Placeholder}' : this.onerror = null;
`
  .trim()
  .replace(/\s+/g, "");
