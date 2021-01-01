

export default class View {
    constructor(anchor) {
        this.fragment = document.createDocumentFragment();

        if (anchor instanceof HTMLDivElement) {
            this.anchor = anchor;
        }
    }

    applyListeners(cloneFragment) {

    }

    render() {
        this.anchor.innerHTML = '';
        const clone = this.fragment.cloneNode(true);
        this.applyListeners(clone);
        return this.anchor.appendChild(clone);
    }
}