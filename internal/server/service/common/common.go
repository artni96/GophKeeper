package common

import commonmodel "github.com/artni96/GophKeeper/internal/server/model/common"

import "fmt"

func EncryptedFieldSetter(pbEntity *commonmodel.PBEncryptedField, entity *commonmodel.EncryptedField, entityName string, notNullFields *[]string) {
	entity.Value = pbEntity.Value.Value
	*notNullFields = append(*notNullFields, fmt.Sprintf("%s_value", entityName))

	entity.Nonce = pbEntity.Nonce.Value
	*notNullFields = append(*notNullFields, fmt.Sprintf("%s_nonce", entityName))

	entity.KeyID = pbEntity.KeyID.Value
	*notNullFields = append(*notNullFields, fmt.Sprintf("%s_key_id", entityName))
}
