class Target < ApplicationRecord
	has_many :drug_targets
	has_many :drugs, through: :drug_targets
	belongs_to :gene
end
