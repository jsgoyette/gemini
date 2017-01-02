package gemini

type Book struct {
	Bids BookEntries
	Asks BookEntries
}

type BookEntries []BookEntry

type BookEntry struct {
	Price  float64 `json:",string"`
	Amount float64 `json:",string"`
}

func (b *BookEntries) Set(price, amount float64) {
	pos := b.findByPrice(price)

	if pos == -1 {
		*b = append(*b, BookEntry{
			Price:  price,
			Amount: amount,
		})
	} else {
		if amount == 0 {
			*b = append((*b)[:pos], (*b)[pos+1:]...)
		} else {
			(*b)[pos].Amount = amount
		}
	}
}

func (b BookEntries) Lowest() BookEntry {

	var lowest float64
	var index int

	if len(b) == 0 {
		return BookEntry{}
	}

	for idx, entry := range b {
		if idx == 0 {
			lowest = entry.Price
			continue
		}
		if entry.Price < lowest {
			lowest = entry.Price
			index = idx
		}
	}

	return b[index]
}

func (b BookEntries) Highest() BookEntry {

	var highest float64
	var index int

	if len(b) == 0 {
		return BookEntry{}
	}

	for idx, entry := range b {
		if idx == 0 {
			highest = entry.Price
			continue
		}
		if entry.Price > highest {
			highest = entry.Price
			index = idx
		}
	}

	return b[index]
}

func (b BookEntries) findByPrice(price float64) int {
	for idx, entry := range b {
		if entry.Price == price {
			return idx
		}
	}
	return -1
}
