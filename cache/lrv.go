package cache

import "sort"

func (c *Cache) LRU() (s Session) {
	var elements []Session

	for	_, s := range c.Elements {
		elements = append(elements, s)	
	} 

	sort.Slice(elements, func(i, j int) bool {
		if elements[i].LastAccess >= elements[j].LastAccess {
			return true
		}	

		return false
	})

	return elements[len(elements) - 1]
}
