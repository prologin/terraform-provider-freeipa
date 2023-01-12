package api

import "errors"

type Group struct {
	CN          []string `json:"cn"`          // Group name
	Description []string `json:"description"` // Group description
}

type GroupList struct {
	CNs []string `json:"memberof_group"` // List of group names
}

func (c *APIClient) GroupAdd(cn string, options JSON) (*Group, error) {
	return apiRequest[Group, string](c, "group_add", options, cn)
}

func (c *APIClient) GroupDel(cn string, options JSON) (*JSON, error) {
	return apiRequest[JSON, []string](c, "group_del", options, cn)
}

func (c *APIClient) GroupMod(cn string, options JSON) (*Group, error) {
	return apiRequest[Group, string](c, "group_mod", options, cn)
}

func (c *APIClient) GroupShow(cn string, options JSON) (*Group, error) {
	return apiRequest[Group, string](c, "group_show", options, cn)
}

func (c *APIClient) GroupFind(criteria string, options JSON) (*[]Group, error) {
	return apiRequest[[]Group, string](c, "group_find", options, criteria)
}

func (c *APIClient) GroupAddMember(cn string, options JSON) (*JSON, error) {
	return apiRequest[JSON, string](c, "group_add_member", options, cn)
}

func (c *APIClient) GroupRemoveMember(cn string, options JSON) (*JSON, error) {
	return apiRequest[JSON, string](c, "group_remove_member", options, cn)
}

func (c *APIClient) GroupAddMemberManager(cn string, options JSON) (*JSON, error) {
	return apiRequest[JSON, string](c, "group_add_member_manager", options, cn)
}

func (c *APIClient) GroupRemoveMemberManager(cn string, options JSON) (*JSON, error) {
	return apiRequest[JSON, string](c, "group_remove_member_manager", options, cn)
}

func (c *APIClient) GetGroups(id string, type_ string) (*GroupList, error) {
	switch type_ {
	case "user":
		return apiRequest[GroupList, string](c, "user_show", nil, id)
	case "group":
		return apiRequest[GroupList, string](c, "group_show", nil, id)
	case "service":
		return apiRequest[GroupList, string](c, "service_show", nil, id)
	}

	return nil, errors.New("invalid type")
}
