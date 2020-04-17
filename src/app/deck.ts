export class Deck{

    constructor(
        public success: boolean = false,
        public deck_id: string = "",
        public shuffled: boolean = false,
        public remaining: number = -1,
        public cards: any[],
        public piles: any
    ){
        this.success = success;
        this.deck_id = deck_id;
        this.shuffled = shuffled;
        this.remaining = remaining;
        this.cards = cards;
        this.piles = piles;
    }
}

export interface DeckGetResponse{
    success: boolean,
    deck_id: string,
    shuffled: boolean,
    remaining: number
  }