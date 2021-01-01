import View from "./View.js";


export default class LoadingView extends View {
    constructor(anchor) {
        super(anchor)

        this.loadingText = document.createElement('h3');
        this.loadingText.innerText = 'Loading...';
        this.loadingText.classList.add('fade');

        this.fragment.appendChild(this.loadingText);
    }
}