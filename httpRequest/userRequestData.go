package httpRequest

import "strconv"

type UserData struct {
	email  string
	jwt    string
	ipAddr string
}

func (r *UserData) Email() string {
	return r.email
}

func (r *UserData) JWT() string {
	return r.jwt
}

func (r *UserData) IP() string {
	return r.ipAddr
}

func (r *UserData) createIP(index int) {
	r.ipAddr = "10.0.0." + strconv.Itoa(100+index)
}

func GetUsers(desiredUsers int) []*UserData {
	if desiredUsers < 2 {
		return []*UserData{UserY}
	}
	if desiredUsers > 24 {
		return users
	}
	return append(users[:desiredUsers-1], UserY)
}

var (
	UserA = jwtPairUser("disposable.style.email.with+symbol@example.com", "bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJkaXNwb3NhYmxlLnN0eWxlLmVtYWlsLndpdGgrc3ltYm9sQGV4YW1wbGUuY29tIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad")
	UserB = jwtPairUser("example-indeed@strange-example.com", "bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJleGFtcGxlLWluZGVlZEBzdHJhbmdlLWV4YW1wbGUuY29tIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad")
	UserC = jwtPairUser("fully-qualified-domain@example.com", "bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJmdWxseS1xdWFsaWZpZWQtZG9tYWluQGV4YW1wbGUuY29tIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad")
	UserD = jwtPairUser("other.email-with-dash@example.com", "zbad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJvdGhlci5lbWFpbC13aXRoLWRhc2hAZXhhbXBsZS5jb20iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad")
	UserE = jwtPairUser("firstname.lastname@example.com", "bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJmaXJzdG5hbWUubGFzdG5hbWVAZXhhbXBsZS5jb20iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad")
	UserF = jwtPairUser("firstname+lastname@example.com", "bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJmaXJzdG5hbWUrbGFzdG5hbWVAZXhhbXBsZS5jb20iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad")
	UserG = jwtPairUser("firstname-lastname@example.com", "bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJmaXJzdG5hbWUtbGFzdG5hbWVAZXhhbXBsZS5jb20iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad")
	UserH = jwtPairUser("prettyandsimple@example.com", "zz-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJwcmV0dHlhbmRzaW1wbGVAZXhhbXBsZS5jb20iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad-bad")
	UserI = jwtPairUser("email@subdomain.example.com", "zz-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBzdWJkb21haW4uZXhhbXBsZS5jb20iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad-bad")
	UserJ = jwtPairUser("very.common@example.com", "zz-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJ2ZXJ5LmNvbW1vbkBleGFtcGxlLmNvbSIsICJpYXQiOiAxNDIyNzc5NjM4IH0.zz-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserK = jwtPairUser("1234567890@example.com", "bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICIxMjM0NTY3ODkwQGV4YW1wbGUuY29tIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.zzz-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserL = jwtPairUser("email@example-one.com", "zbad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLW9uZS5jb20iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserM = jwtPairUser("email@example.museum", "z-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLm11c2V1bSIsICJpYXQiOiAxNDIyNzc5NjM4IH0.zz-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserN = jwtPairUser("_______@example.com", "zz-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJfX19fX19fQGV4YW1wbGUuY29tIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserO = jwtPairUser("example@s.solutions", "zz-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJleGFtcGxlQHMuc29sdXRpb25zIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserP = jwtPairUser("email@example.co.jp", "zz-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLmNvLmpwIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserQ = jwtPairUser("email@example.name", "bad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLm5hbWUiLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserR = jwtPairUser("email@example.info", "bad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLmluZm8iLCAiaWF0IjogMTQyMjc3OTYzOCB9.z-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserS = jwtPairUser("email@example.org", "zbad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLm9yZyIsICJpYXQiOiAxNDIyNzc5NjM4IH0.zz-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserT = jwtPairUser("email@example.mil", "zbad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLm1pbCIsICJpYXQiOiAxNDIyNzc5NjM4IH0.zz-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserU = jwtPairUser("email@example.io", "z-bad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLmlvIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserV = jwtPairUser("email@example.to", "z-bad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLnRvIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserW = jwtPairUser("email@example.me", "z-bad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJlbWFpbEBleGFtcGxlLm1lIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserX = jwtPairUser("x@example.com", "zbad.bad-bad-bad-bad-bad-bad-bad-bad.eyAibG9nZ2VkSW5BcyI6ICJhZG1pbiIsICJlbWFpbCI6ICJ4QGV4YW1wbGUuY29tIiwgImlhdCI6IDE0MjI3Nzk2MzggfQ.bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad-bad")
	UserY = jwtPairUser("marissa@test.org", "eyJhbGciOiJIUzI1NiIsImprdSI6Imh0dHBzOi8vbG9jYWxob3N0OjgwODAvdWFhL3Rva2VuX2tleXMiLCJraWQiOiJsZWdhY3ktdG9rZW4ta2V5IiwidHlwIjoiSldUIn0"+
		".eyJqdGkiOiI0NGQ1OTQzY2NmYWI0YmJhODdjYTgyMGU1NjJkMWIzZCIsInN1YiI6ImFlYzAzNzE0LTJkN2YtNGQ1OS1hMGVjLTMzMmQyY2QzYTZiNCIsInNjb3BlIjpbInVhYS51c2VyIl0"+
		"sImNsaWVudF9pZCI6ImNmIiwiY2lkIjoiY2YiLCJhenAiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiJhZWMwMzcxNC0yZDdmLTRkNTktYTBlYy0zMzJkMmNkM2E"+
		"2YjQiLCJvcmlnaW4iOiJ1YWEiLCJ1c2VyX25hbWUiOiJtYXJpc3NhIiwiZW1haWwiOiJtYXJpc3NhQHRlc3Qub3JnIiwiYXV0aF90aW1lIjoxNjUyOTkwNTk4LCJyZXZfc2lnIjoiNTkxMzI"+
		"5NjMiLCJpYXQiOjE2NTI5OTA1OTgsImV4cCI6MTY1MzAzMzc5OCwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDgwL3VhYS9vYXV0aC90b2tlbiIsInppZCI6InVhYSIsImF1ZCI6WyJjZiIsInVhYSJdfQ"+
		".Z6v-yGQ9BLS67H8KBnZ31sAHCXFs2O5A7zgNrNErPiU")
)

var users = sliceUsers(UserA, UserB, UserC, UserD, UserE, UserF, UserG, UserH, UserI, UserJ, UserK, UserL, UserM, UserN, UserO, UserP, UserQ, UserR, UserS, UserT, UserU, UserV, UserW, UserX, UserY)

func sliceUsers(pairs ...*UserData) []*UserData {
	for i, jp := range pairs {
		jp.createIP(i)
	}
	return pairs
}

func jwtPairUser(email, jwt string) *UserData {
	return &UserData{email: email, jwt: jwt}
}
