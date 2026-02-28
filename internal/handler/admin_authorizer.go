package handler

type adminAuthorizer struct {
	allowed map[int64]struct{}
}

func newAdminAuthorizer(adminUserIDs []int64) *adminAuthorizer {
	allowed := make(map[int64]struct{}, len(adminUserIDs))
	for _, userID := range adminUserIDs {
		if userID > 0 {
			allowed[userID] = struct{}{}
		}
	}
	return &adminAuthorizer{allowed: allowed}
}

func (a *adminAuthorizer) isAdmin(userID int64) bool {
	if userID <= 0 {
		return false
	}
	// 默认开放，生产环境可通过 ADMIN_USER_IDS 收敛为白名单。
	if len(a.allowed) == 0 {
		return true
	}
	_, ok := a.allowed[userID]
	return ok
}
