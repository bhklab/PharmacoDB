class Gene < ApplicationRecord
	has_many :gene_drugs
	has_many :drugs, through: :gene_drugs
	has_one :target
end
