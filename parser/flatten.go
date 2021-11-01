package parser

import (
	"github.com/hashicorp/hcl/hcl/ast"
)

// MEMO : https://github.com/hashicorp/hcl/tree/hcl1/json/parser가 원본. 코드 수정해서 사용하기 위해 카피 함.

// MEMO : 임시 처리. 추후 수정 예정.
var independentBlocksKeys = []string{"\"egress\"", "\"ingress\"", "\"filter\""}

func stringArrayContains(items []string, item string) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}

	return false
}

// flattenObjects takes an AST node, walks it, and flattens
func flattenObjects(node ast.Node) {
	ast.Walk(node, func(n ast.Node) (ast.Node, bool) {
		// We only care about lists, because this is what we modify
		list, ok := n.(*ast.ObjectList)
		if !ok {
			return n, true
		}

		// Rebuild the item list
		items := make([]*ast.ObjectItem, 0, len(list.Items))
		frontier := make([]*ast.ObjectItem, len(list.Items))
		copy(frontier, list.Items)
		for len(frontier) > 0 {
			// Pop the current item
			n := len(frontier)
			item := frontier[n-1]
			frontier = frontier[:n-1]
			switch v := item.Val.(type) {
			case *ast.ObjectType:
				items, frontier = flattenObjectType(v, item, items, frontier)
			case *ast.ListType:
				items, frontier = flattenListType(v, item, items, frontier)
			default:
				items = append(items, item)
			}
		}
		// Reverse the list since the frontier model runs things backwards
		for i := len(items)/2 - 1; i >= 0; i-- {
			opp := len(items) - 1 - i
			items[i], items[opp] = items[opp], items[i]
		}

		// Done! Set the original items
		list.Items = items
		return n, true
	})
}

func flattenListType(
	ot *ast.ListType,
	item *ast.ObjectItem,
	items []*ast.ObjectItem,
	frontier []*ast.ObjectItem) ([]*ast.ObjectItem, []*ast.ObjectItem) {
	// If the list is empty, keep the original list
	if len(ot.List) == 0 {
		items = append(items, item)
		return items, frontier
	}
	// MEMO : egress, ingress 속성의 경우 json상에선 배열형태이지만, 역변환 tf내용은 각각의 block으로 들어가야 함.
	// MEMO : 그 땐 배열화해주는 부분 타지 않도록 조건문 추가함
	if !stringArrayContains(independentBlocksKeys, item.Keys[0].Token.Text) {
		// All the elements of this object must also be objects!
		for _, subitem := range ot.List {
			// MEMO : ast.ObjectType으로 하면 object들로 이루어진 배열일 때 여기에서 append안돼서 정상변환 안돼서 이렇게 수정함.
			if _, ok := subitem.(*ast.ListType); !ok {
				items = append(items, item)
				return items, frontier
			}
		}
	}

	// Great! We have a match go through all the items and flatten
	for _, elem := range ot.List {
		// Add it to the frontier so that we can recurse
		frontier = append(frontier, &ast.ObjectItem{
			Keys:        item.Keys,
			Assign:      item.Assign,
			Val:         elem,
			LeadComment: item.LeadComment,
			LineComment: item.LineComment,
		})
	}

	return items, frontier
}

func flattenObjectType(
	ot *ast.ObjectType,
	item *ast.ObjectItem,
	items []*ast.ObjectItem,
	frontier []*ast.ObjectItem) ([]*ast.ObjectItem, []*ast.ObjectItem) {

	// If the list has no items we do not have to flatten anything
	if ot.List.Items == nil {
		items = append(items, item)
		return items, frontier
	}

	for _, subitem := range ot.List.Items {
		if _, ok := subitem.Val.(*ast.ObjectType); !ok {
			// MEMO : provider 키워드로 하위에 들어가있는 정의가 배열형태일 땐 각각을 provider '레이블' {} 블럭으로 생성해줘야 돼서 그 예외처리 부분 추가함.
			if item.Keys[0].Token.Text == "\"provider\"" {
				if !contains(items, item) && len(item.Keys) > 1 {
					items = append(items, item)
				}
				if _, ok := subitem.Val.(*ast.ListType); ok {
					keys := make([]*ast.ObjectKey, len(item.Keys)+len(subitem.Keys))
					copy(keys, item.Keys)
					copy(keys[len(item.Keys):], subitem.Keys)
					frontier = appenChildAsEachObject(item, keys, subitem.Val, frontier)
				}
			} else {
				if !contains(items, item) {
					items = append(items, item)
				}
			}
			return items, frontier
		} else {
			// MEMO : variable의 경우 하위 속성들이 object로만 이루어져있어도 역변환 결과물이 variable A B {} 형태로 나오면 안됨.
			// MEMO : 무조건 variable A {} 형태로 나오게 하기 위해 아래 구문 추가함.
			if len(item.Keys) > 1 && item.Keys[0].Token.Text == "\"variable\"" {
				if !contains(items, item) {
					items = append(items, item)
				}
				return items, frontier
			}
		}
	}

	for _, subitem := range ot.List.Items {
		{
			keys := make([]*ast.ObjectKey, len(item.Keys)+len(subitem.Keys))
			copy(keys, item.Keys)
			copy(keys[len(item.Keys):], subitem.Keys)
			frontier = appenChildAsEachObject(item, keys, subitem.Val, frontier)
		}
	}
	return items, frontier
}

func contains(items []*ast.ObjectItem, item *ast.ObjectItem) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}

	return false
}

func appenChildAsEachObject(item *ast.ObjectItem, keys []*ast.ObjectKey, n interface{}, frontier []*ast.ObjectItem) []*ast.ObjectItem {
	// var frontier []*ast.ObjectItem
	switch t := n.(type) {
	case *ast.ListType:
		for _, i := range t.List {
			frontier = append(frontier, &ast.ObjectItem{
				Keys:        keys,
				Assign:      item.Assign,
				Val:         i,
				LeadComment: item.LeadComment,
				LineComment: item.LineComment,
			})
		}
	case *ast.ObjectType:
		frontier = append(frontier, &ast.ObjectItem{
			Keys:        keys,
			Assign:      item.Assign,
			Val:         t,
			LeadComment: item.LeadComment,
			LineComment: item.LineComment,
		})
	}
	return frontier
}
